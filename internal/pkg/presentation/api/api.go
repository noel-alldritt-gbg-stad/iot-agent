package api

import (
	"io/ioutil"
	"net/http"

	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/logging"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("iot-agent/api")

type API interface {
	Start(port string) error
	health(w http.ResponseWriter, r *http.Request)
}

type api struct {
	log zerolog.Logger
	r   chi.Router
	app iotagent.IoTAgent
}

func NewApi(logger zerolog.Logger, r chi.Router, app iotagent.IoTAgent) API {
	a := newAPI(logger, r, app)

	return a
}

func newAPI(logger zerolog.Logger, r chi.Router, app iotagent.IoTAgent) *api {
	a := &api{
		log: logger,
		r:   r,
		app: app,
	}

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	serviceName := "iot-agent"

	r.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(r)))

	r.Get("/health", a.health)
	r.Post("/newmsg", a.incomingMsg)

	return a
}

func (a *api) Start(port string) error {
	a.log.Info().Str("port", port).Msg("starting to listen for connections")

	return http.ListenAndServe(":"+port, a.r)
}

func (a *api) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (a *api) incomingMsg(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx, span := tracer.Start(r.Context(), "newmsg")
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	log := a.log

	traceID := span.SpanContext().TraceID()
	if traceID.IsValid() {
		log = log.With().Str("traceID", traceID.String()).Logger()
	}

	ctx = logging.NewContextWithLogger(ctx, log)

	msg, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	log.Info().Msg("attempting to process message")
	err = a.app.MessageReceived(ctx, msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		log.Error().Err(err).Msg("failed to handle message")
		return
	}

	w.WriteHeader(http.StatusCreated)
}

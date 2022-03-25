package api

import (
	"compress/flate"
	"io/ioutil"
	"net/http"

	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

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

	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	r.Use(compressor.Handler)

	r.Use(httplog.RequestLogger(httplog.NewLogger("iot-agent", httplog.Options{
		JSON: true,
	})))

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
	msg, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	err := a.app.MessageReceived(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

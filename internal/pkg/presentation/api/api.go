package api

import (
	"compress/flate"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

type API interface {
	Start(port string) error
	health(w http.ResponseWriter, r *http.Request)
}

type iotAgentApi struct {
	r chi.Router
}

func NewApi(r chi.Router) API {

	a := newIotAgentApi(r)

	return a
}

func newIotAgentApi(r chi.Router) *iotAgentApi {
	a := &iotAgentApi{
		r: r,
	}

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	r.Use(compressor.Handler)

	logger := httplog.NewLogger("iot-agent", httplog.Options{
		JSON: true,
	})
	r.Use(httplog.RequestLogger(logger))

	r.Get("/health", a.health)

	return a
}

func (a *iotAgentApi) Start(port string) error {
	log.Info().Str("port", port).Msg("starting to listen for connections")

	return http.ListenAndServe(":"+port, a.r)
}

func (a *iotAgentApi) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

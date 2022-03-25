package main

import (
	"os"
	"strings"

	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	serviceName := "iot-agent"

	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()
	logger.Info().Msg("starting up ...")

	app := SetupIoTAgent()

	mqttConfig, _ := mqtt.NewConfigFromEnvironment()
	mqttClient, err := mqtt.NewClient(logger, mqttConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create mqtt client")
	}
	mqttClient.Start()
	defer mqttClient.Stop()

	SetupAndRunApi(logger, app)
}

func SetupIoTAgent() iotagent.IoTAgent {
	dmc := domain.NewDeviceManagementClient()
	event := events.NewEventPublisher()

	return iotagent.NewIoTAgent(dmc, event)
}

func SetupAndRunApi(logger zerolog.Logger, app iotagent.IoTAgent) {
	r := chi.NewRouter()

	a := api.NewApi(logger, r, app)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8880"
	}

	a.Start(port)
}

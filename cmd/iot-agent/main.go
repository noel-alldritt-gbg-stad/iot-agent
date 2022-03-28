package main

import (
	"os"
	"runtime/debug"
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
	logger := newLogger("iot-agent")
	logger.Info().Msg("starting up ...")

	app := SetupIoTAgent()

	apiPort := os.Getenv("SERVICE_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	mqttConfig, err := mqtt.NewConfigFromEnvironment()
	if err != nil {
		logger.Fatal().Err(err).Msg("mqtt configuration error")
	}

	mqttClient, err := mqtt.NewClient(logger, mqttConfig, apiPort)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create mqtt client")
	}
	mqttClient.Start()
	defer mqttClient.Stop()

	SetupAndRunApi(logger, app, apiPort)
}

func newLogger(serviceName string) zerolog.Logger {
	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		buildSettings := buildInfo.Settings
		infoMap := map[string]string{}
		for _, s := range buildSettings {
			infoMap[s.Key] = s.Value
		}

		sha := infoMap["vcs.revision"]
		if infoMap["vcs.modified"] == "true" {
			sha += "+"
		}

		logger = logger.With().Str("version", sha).Logger()
	} else {
		logger.Error().Msg("failed to extract build information")
	}

	return logger
}

func SetupIoTAgent() iotagent.IoTAgent {
	dmc := domain.NewDeviceManagementClient()
	event := events.NewEventPublisher()

	return iotagent.NewIoTAgent(dmc, event)
}

func SetupAndRunApi(logger zerolog.Logger, app iotagent.IoTAgent, port string) {
	r := chi.NewRouter()

	a := api.NewApi(logger, r, app)

	a.Start(port)
}

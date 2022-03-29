package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/tracing"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	serviceName := "iot-agent"
	serviceVersion := version()

	logger := newLogger(serviceName, serviceVersion)
	logger.Info().Msg("starting up ...")

	ctx := context.Background()

	cleanup, err := tracing.Init(ctx, logger, serviceName, serviceVersion)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init tracing")
	}
	defer cleanup()

	app := SetupIoTAgent(logger)

	apiPort := os.Getenv("SERVICE_PORT")
	if apiPort == "" {
		apiPort = "8080"
	}

	mqttConfig, err := mqtt.NewConfigFromEnvironment()
	if err != nil {
		logger.Fatal().Err(err).Msg("mqtt configuration error")
	}

	forwardingEndpoint := fmt.Sprintf("http://127.0.0.1:%s/newmsg", apiPort)
	mqttClient, err := mqtt.NewClient(logger, mqttConfig, forwardingEndpoint)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create mqtt client")
	}
	mqttClient.Start()
	defer mqttClient.Stop()

	SetupAndRunApi(logger, app, apiPort)
}

func newLogger(serviceName, serviceVersion string) zerolog.Logger {
	logger := log.With().Str("service", strings.ToLower(serviceName)).Str("version", serviceVersion).Logger()
	return logger
}

func version() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	buildSettings := buildInfo.Settings
	infoMap := map[string]string{}
	for _, s := range buildSettings {
		infoMap[s.Key] = s.Value
	}

	sha := infoMap["vcs.revision"]
	if infoMap["vcs.modified"] == "true" {
		sha += "+"
	}

	return sha
}

func SetupIoTAgent(logger zerolog.Logger) iotagent.IoTAgent {
	dmcUrl := os.Getenv("DMC_URL")
	if dmcUrl == "" {
		logger.Fatal().Msgf("DMC_URL must be set")
	}
	dmc := domain.NewDeviceManagementClient(dmcUrl, logger)
	event := events.NewEventPublisher()

	return iotagent.NewIoTAgent(dmc, event, logger)
}

func SetupAndRunApi(logger zerolog.Logger, app iotagent.IoTAgent, port string) {
	r := chi.NewRouter()

	a := api.NewApi(logger, r, app)

	a.Start(port)
}

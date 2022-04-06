package main

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/logging"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/tracing"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

func main() {
	serviceName := "iot-agent"
	serviceVersion := version()

	ctx, logger := logging.NewLogger(context.Background(), serviceName, serviceVersion)
	logger.Info().Msg("starting up ...")

	cleanup, err := tracing.Init(ctx, logger, serviceName, serviceVersion)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init tracing")
	}
	defer cleanup()

	app := SetupIoTAgent(serviceName, logger)

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

func SetupIoTAgent(serviceName string, logger zerolog.Logger) iotagent.IoTAgent {
	dmcUrl := os.Getenv("DEV_MGMT_URL")
	dmc := domain.NewDeviceManagementClient(dmcUrl)
	event := events.NewEventSender(serviceName, logger)
	event.Start()

	return iotagent.NewIoTAgent(dmc, event, logger)
}

func SetupAndRunApi(logger zerolog.Logger, app iotagent.IoTAgent, port string) {
	r := chi.NewRouter()

	a := api.NewApi(logger, r, app)

	a.Start(port)
}

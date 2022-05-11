package main

import (
	"context"

	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/diwise/service-chassis/pkg/infrastructure/buildinfo"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

const serviceName string = "iot-agent"

func main() {

	serviceVersion := buildinfo.SourceVersion()
	_, logger, cleanup := o11y.Init(context.Background(), serviceName, serviceVersion)
	defer cleanup()

	forwardingEndpoint := env.GetVariableOrDie(logger, "MSG_FWD_ENDPOINT", "endpoint that incoming packages should be forwarded to")
	deviceMgmtClientURL := env.GetVariableOrDie(logger, "DEV_MGMT_URL", "device management client URL")

	apiPort := env.GetVariableOrDefault(logger, "SERVICE_PORT", "8080")

	app := SetupIoTAgent(logger, serviceName, deviceMgmtClientURL)

	mqttConfig, err := mqtt.NewConfigFromEnvironment()
	if err != nil {
		logger.Fatal().Err(err).Msg("mqtt configuration error")
	}

	mqttClient, err := mqtt.NewClient(logger, mqttConfig, forwardingEndpoint)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create mqtt client")
	}
	mqttClient.Start()
	defer mqttClient.Stop()

	SetupAndRunApi(logger, app, apiPort)
}

func SetupIoTAgent(logger zerolog.Logger, serviceName, deviceMgmtClientURL string) iotagent.IoTAgent {
	dmc := domain.NewDeviceManagementClient(deviceMgmtClientURL)
	event := events.NewEventSender(serviceName, logger)
	event.Start()

	return iotagent.NewIoTAgent(dmc, event)
}

func SetupAndRunApi(logger zerolog.Logger, app iotagent.IoTAgent, port string) {
	r := chi.NewRouter()

	a := api.NewApi(logger, r, app)

	metrics.AddHandlers(r)

	a.Start(port)
}

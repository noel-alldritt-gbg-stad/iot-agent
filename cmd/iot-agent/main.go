package main

import (
	"os"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/iotagent"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	app := SetupIoTAgent()

	mqttConfig, _ := mqtt.NewConfigFromEnvironment()
	mqttClient, err := mqtt.NewClient(mqttConfig)
	if err != nil {
		panic("failed to create mqtt client: " + err.Error()) // TODO: Use proper logging (will be handled in its own commit)
	}
	mqttClient.Start()
	defer mqttClient.Stop()

	SetupAndRunApi(app)
}

func SetupIoTAgent() iotagent.IoTAgent {
	dmc := domain.NewDeviceManagementClient()
	cr := conversion.NewConverterRegistry()
	event := events.NewEventPublisher()
	mp := messageprocessor.NewMessageReceivedProcessor(dmc, cr, event)

	return iotagent.NewIoTAgent(mp)
}

func SetupAndRunApi(app iotagent.IoTAgent) {
	r := chi.NewRouter()

	a := api.NewApi(r, app)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8880"
	}

	a.Start(port)
}

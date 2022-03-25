package main

import (
	"os"

	"github.com/diwise/iot-agent/internal/pkg/application"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/services/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	app := SetupIoTAgent()

	go SetupAndRunApi(app)

	mqtt.SetupAndRunMQTT()
}

func SetupIoTAgent() application.IoTAgent {
	dmc := domain.NewDeviceManagementClient()
	cr := application.NewConverterRegistry()
	event := application.NewEventPublisher()
	mp := application.MessageReceivedProcessor(dmc, cr, event)

	return application.NewIoTAgent(mp)
}

func SetupAndRunApi(app application.IoTAgent) {
	r := chi.NewRouter()

	a := api.NewApi(r, app)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8880"
	}

	a.Start(port)
}

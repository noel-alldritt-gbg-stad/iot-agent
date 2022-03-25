package main

import (
	"os"

	"github.com/diwise/iot-agent/internal/pkg/application"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
)

func main() {
	mp := application.NewMessageProcessor()
	app := application.NewIoTAgent(mp)

	go SetupAndRunApi(app)

	mqtt.SetupAndRunMQTT()
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

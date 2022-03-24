package main

import (
	"os"

	"github.com/diwise/iot-agent/internal/pkg/infrastructure/mqtt"
	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	"github.com/go-chi/chi/v5"
)

func main() {

	go mqtt.SetupAndRunMQTT()

	SetupAndRunApi()
}

func SetupAndRunApi() {
	r := chi.NewRouter()

	a := api.NewApi(r)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8880"
	}

	a.Start(port)
}

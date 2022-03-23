package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/diwise/iot-agent/internal/pkg/presentation/api"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func main() {

	go SetupAndRunMQTT()

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

func SetupAndRunMQTT() {
	options := mqtt.NewClientOptions()
	// broker IP and port

	connectionString := fmt.Sprintf("tls://%s:8883", os.Getenv("MQTT_HOST"))
	options.AddBroker(connectionString)

	options.Username = os.Getenv("MQTT_USER")
	options.Password = os.Getenv("MQTT_PASSWORD")

	options.SetClientID("diwise/iot-agent" + uuid.NewString())
	options.SetDefaultPublishHandler(MessageHandler)

	options.OnConnect = func(c mqtt.Client) {
		fmt.Println("connected!")
		c.Subscribe("application/53/device/#", 0, nil)
	}

	options.OnConnectionLost = func(c mqtt.Client, err error) {
		panic(fmt.Sprintf("connection lost: %s\n", err.Error()))
	}

	options.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		time.Sleep(1 * time.Second)
	}
}

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	fmt.Printf("received payload %s", string(payload))
	msg.Ack()
}

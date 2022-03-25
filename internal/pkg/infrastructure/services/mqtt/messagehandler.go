package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func MessageHandler(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()
	fmt.Printf("received payload %s", string(payload))
	msg.Ack()
}

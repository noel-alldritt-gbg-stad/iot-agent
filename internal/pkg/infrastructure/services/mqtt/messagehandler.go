package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func NewMessageHandler(logger zerolog.Logger) func(mqtt.Client, mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		logger.Info().Msgf("received payload: %s", string(payload))
		msg.Ack()
	}
}

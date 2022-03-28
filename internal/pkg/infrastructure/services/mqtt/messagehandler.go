package mqtt

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

func NewMessageHandler(logger zerolog.Logger, apiPort string) func(mqtt.Client, mqtt.Message) {
	messageReceiver := fmt.Sprintf("http://127.0.0.1:%s/newmsg", apiPort)

	return func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()

		httpClient := http.Client{
			Transport: nil, //otelhttp.NewTransport(http.DefaultTransport),
		}

		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, messageReceiver, bytes.NewBuffer(payload))
		if err != nil {
			return
		}

		req.Header.Add("Content-Type", "application/json")

		logger.Info().Msgf("forwarding received payload to %s", messageReceiver)
		_, err = httpClient.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("failed to forward message")
		}

		msg.Ack()
	}
}

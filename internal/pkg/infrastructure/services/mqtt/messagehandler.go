package mqtt

import (
	"bytes"
	"context"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("iot-agent/mqtt/message-handler")

func NewMessageHandler(logger zerolog.Logger, forwardingEndpoint string) func(mqtt.Client, mqtt.Message) {

	return func(client mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()

		httpClient := http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}

		var err error

		ctx, span := tracer.Start(context.Background(), "forward-message")
		defer func() {
			if err != nil {
				span.RecordError(err)
			}
			span.End()
		}()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, forwardingEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			logger.Error().Err(err).Msg("failed to create http request")
			return
		}

		req.Header.Add("Content-Type", "application/json")

		logger.Info().Msgf("forwarding received payload to %s", forwardingEndpoint)
		_, err = httpClient.Do(req)
		if err != nil {
			logger.Error().Err(err).Msg("failed to forward message")
		}

		msg.Ack()
	}
}

package mqtt

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
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
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

		log := logger

		traceID := span.SpanContext().TraceID()
		if traceID.IsValid() {
			log = logger.With().Str("traceID", traceID.String()).Logger()
		}

		log.Info().Msgf("received payload %s from topic %s", string(payload), msg.Topic())

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, forwardingEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			log.Error().Err(err).Msg("failed to create http request")
			return
		}

		req.Header.Add("Content-Type", "application/json")

		log.Info().Msgf("forwarding received payload to %s", forwardingEndpoint)
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("forwarding request failed")
		} else if resp.StatusCode != http.StatusCreated {
			err = fmt.Errorf("unexpected response code %d", resp.StatusCode)
			log.Error().Err(err).Msg("failed to forward message")
		}

		msg.Ack()
	}
}

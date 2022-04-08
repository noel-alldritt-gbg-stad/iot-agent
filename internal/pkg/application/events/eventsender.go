package events

import (
	"context"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/logging"
	"github.com/rs/zerolog"

	"github.com/diwise/messaging-golang/pkg/messaging"
)

type EventSender interface {
	Start() error
	Send(ctx context.Context, msg conversion.InternalMessage) error
	Stop() error
}

type eventSender struct {
	rmqConfig    messaging.Config
	rmqMessenger messaging.MsgContext
	started      bool
}

func NewEventSender(serviceName string, logger zerolog.Logger) EventSender {
	sender := &eventSender{
		rmqConfig: messaging.LoadConfiguration(serviceName, logger),
	}

	return sender
}

func (e *eventSender) Send(ctx context.Context, msg conversion.InternalMessage) error {
	log := logging.GetFromContext(ctx)

	if !e.started {
		err := fmt.Errorf("attempt to send before start")
		log.Error().Err(err).Msg("send failed")
		return err
	}

	log.Info().Msg("sending command to iot-core queue")
	return e.rmqMessenger.SendCommandTo(ctx, msg, "iot-core")
}

func (e *eventSender) Start() error {
	var err error
	e.rmqMessenger, err = messaging.Initialize(e.rmqConfig)
	if err == nil {
		e.started = true
	}
	return err
}

func (e *eventSender) Stop() error {
	e.rmqMessenger.Close()
	return nil
}

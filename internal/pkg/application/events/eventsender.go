package events

import (
	"context"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/rs/zerolog"

	"github.com/diwise/messaging-golang/pkg/messaging"
)

type EventSender interface {
	Start() error
	Send(ctx context.Context, msg conversion.InternalMessage) error
	Stop() error
}

type eventSender struct {
	logger       zerolog.Logger
	rmqConfig    messaging.Config
	rmqMessenger messaging.MsgContext
	started      bool
}

func NewEventSender(serviceName string, logger zerolog.Logger) EventSender {
	sender := &eventSender{
		logger:    logger,
		rmqConfig: messaging.LoadConfiguration(serviceName, logger),
	}

	return sender
}

func (e *eventSender) Send(ctx context.Context, msg conversion.InternalMessage) error {
	if !e.started {
		err := fmt.Errorf("attempt to send before start")
		e.logger.Error().Err(err).Msg("send failed")
		return err
	}

	e.logger.Info().Msg("publishing message on topic " + msg.TopicName())
	return e.rmqMessenger.SendCommandTo(ctx, msg, "msg.rcvd")
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

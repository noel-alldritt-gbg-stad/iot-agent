package events

import (
	"context"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/rs/zerolog"

	"github.com/diwise/messaging-golang/pkg/messaging"
)

type EventPublisher interface {
	Start() error
	Publish(ctx context.Context, msg conversion.InternalMessage) error
	Stop() error
}

type eventPublisher struct {
	logger       zerolog.Logger
	rmqConfig    messaging.Config
	rmqMessenger messaging.MsgContext
	started      bool
}

func NewEventPublisher(serviceName string, logger zerolog.Logger) EventPublisher {
	publisher := &eventPublisher{
		logger:    logger,
		rmqConfig: messaging.LoadConfiguration(serviceName, logger),
	}

	return publisher
}

//places a converted message on rabbit...
func (e *eventPublisher) Publish(ctx context.Context, msg conversion.InternalMessage) error {
	if !e.started {
		err := fmt.Errorf("attempt to publish before start")
		e.logger.Error().Err(err).Msg("publish failed")
		return err
	}

	return nil
}

func (e *eventPublisher) Start() error {
	var err error
	e.rmqMessenger, err = messaging.Initialize(e.rmqConfig)
	if err == nil {
		e.started = true
	}
	return err
}

func (e *eventPublisher) Stop() error {
	e.rmqMessenger.Close()
	return nil
}

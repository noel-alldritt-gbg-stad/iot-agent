package events

import (
	"context"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/rs/zerolog"

	"github.com/diwise/messaging-golang/pkg/messaging"
)

type EventPublisher interface {
	Start() error
	Publish(ctx context.Context, msg conversion.InternalMessageFormat) error
	Stop() error
}

type eventPublisher struct {
	logger       zerolog.Logger
	rmqConfig    messaging.Config
	rmqMessenger messaging.MsgContext
}

func NewEventPublisher(serviceName string, logger zerolog.Logger) EventPublisher {
	publisher := &eventPublisher{
		logger:    logger,
		rmqConfig: messaging.LoadConfiguration(serviceName, logger),
	}

	return publisher
}

//places a converted message on rabbit...
func (e *eventPublisher) Publish(ctx context.Context, msg conversion.InternalMessageFormat) error {
	e.logger.Error().Msg("publishing to queue is not yet implemented.")
	return nil
}

func (e *eventPublisher) Start() error {
	var err error
	e.rmqMessenger, err = messaging.Initialize(e.rmqConfig)
	return err
}

func (e *eventPublisher) Stop() error {
	e.rmqMessenger.Close()
	return nil
}

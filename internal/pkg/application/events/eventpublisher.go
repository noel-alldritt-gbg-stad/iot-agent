package events

import (
	"context"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
)

type EventPublisher interface {
	Publish(ctx context.Context, msg conversion.InternalMessageFormat) error
}

type eventPublisher struct {
}

func NewEventPublisher() EventPublisher {
	return &eventPublisher{}
}

//places a converted message on rabbit...
func (e *eventPublisher) Publish(ctx context.Context, msg conversion.InternalMessageFormat) error {

	return nil
}

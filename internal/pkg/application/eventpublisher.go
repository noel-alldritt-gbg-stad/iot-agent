package application

type EventPublisher interface {
}

type eventPublisher struct {
}

func NewEventPublisher() EventPublisher {
	event := &eventPublisher{}

	return event
}

//places a converted message on rabbit...

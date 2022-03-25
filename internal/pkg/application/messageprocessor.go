package application

import "fmt"

type MessageProcessor interface {
	ProcessMessage(msg []byte) error
}

// hantera kö av msgs, skicka till converter registry

type msgProcessor struct {
}

func NewMessageProcessor() MessageProcessor {
	mp := &msgProcessor{}

	return mp
}

func (mp *msgProcessor) ProcessMessage(msg []byte) error {
	return fmt.Errorf("not implemented yet")
}

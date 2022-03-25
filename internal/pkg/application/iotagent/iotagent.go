package iotagent

import (
	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	"github.com/diwise/iot-agent/internal/pkg/domain"
)

type IoTAgent interface {
	MessageReceived(msg []byte) error
}

type iotAgent struct {
	mp messageprocessor.MessageProcessor
}

func NewIoTAgent(dmc domain.DeviceManagementClient, eventPub events.EventPublisher) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub)

	return &iotAgent{
		mp: msgprcs,
	}
}

// this function is likely to be renamed
func (a *iotAgent) MessageReceived(msg []byte) error {
	err := a.mp.ProcessMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

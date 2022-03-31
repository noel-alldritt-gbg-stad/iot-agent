package iotagent

import (
	"context"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/rs/zerolog"
)

type IoTAgent interface {
	MessageReceived(ctx context.Context, msg []byte) error
}

type iotAgent struct {
	mp messageprocessor.MessageProcessor
}

func NewIoTAgent(dmc domain.DeviceManagementClient, eventPub events.EventSender, log zerolog.Logger) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub, log)

	return &iotAgent{
		mp: msgprcs,
	}
}

func (a *iotAgent) MessageReceived(ctx context.Context, msg []byte) error {
	err := a.mp.ProcessMessage(ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

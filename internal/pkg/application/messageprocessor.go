package application

import (
	"encoding/json"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/domain"
)

type MessageProcessor interface {
	ProcessMessage(msg []byte) error
}

// hantera kö av msgs, skicka till converter registry

type msgProcessor struct {
	dmc    domain.DeviceManagementClient
	conReg ConverterRegistry
	event  EventPublisher
}

func NewMessageProcessor(dmc domain.DeviceManagementClient, conReg ConverterRegistry, event EventPublisher) MessageProcessor {
	mp := &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
	}

	return mp
}

func (mp *msgProcessor) ProcessMessage(msg []byte) error {
	// extract and send devEUI to devicemanagementclient

	dm := DeviceMessage{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return err
	}

	mp.dmc.FindDeviceFromDevEUI()

	// response with internal id, type and format gets passed to Converter registry

	mp.conReg.Designate()

	// converter registry returns msg payload in internal format

	// converted message gets sent to event publisher

	return fmt.Errorf("not implemented yet")
}

type DeviceMessage struct {
	DevEUI string `json:"devEUI"`
}

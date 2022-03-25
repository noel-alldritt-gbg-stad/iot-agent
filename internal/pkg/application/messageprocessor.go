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

func MessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg ConverterRegistry, event EventPublisher) MessageProcessor {
	mp := &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
	}

	return mp
}

func (mp *msgProcessor) ProcessMessage(msg []byte) error {
	// extract and send devEUI to devicemanagementclient
	// format is from mqtt, not device management client

	dm := DeviceMessage{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return err
	}

	mp.dmc.FindDeviceFromDevEUI()

	// response with internal id, type and gets passed to Converter registry

	mp.conReg.Designate()

	// converter registry returns the correct converter

	// msg converter converts msg payload to internal format and returns it

	// converted message gets sent to event publisher

	return fmt.Errorf("not implemented yet")
}

type DeviceMessage struct {
	DevEUI string `json:"devEUI"`
}

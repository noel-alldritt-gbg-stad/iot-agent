package messageprocessor

import (
	"encoding/json"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
)

type MessageProcessor interface {
	ProcessMessage(msg []byte) error
}

// hantera kö av msgs, skicka till converter registry

type msgProcessor struct {
	dmc    domain.DeviceManagementClient
	conReg conversion.ConverterRegistry
	event  events.EventPublisher
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventPublisher) MessageProcessor {
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
	if err == nil {
		result, err := mp.dmc.FindDeviceFromDevEUI(dm.DevEUI)
		if err == nil {
			messageConverter := mp.conReg.DesignateConverters(result)
			if len(messageConverter) == 0 {
				return fmt.Errorf("no matching converters for device")
			}

			for _, mc := range messageConverter {
				// msg converter converts msg payload to internal format and returns it
				payload, err := mc.ConvertPayload(msg)
				if err == nil {
					mp.event.Publish(payload)
					// converted message gets sent to event publisher
				}

			}
		}
	}

	// response with internal id, type and gets passed to Converter registry
	// converter registry returns the correct converters

	return err
}

type DeviceMessage struct {
	DevEUI string `json:"devEUI"`
}

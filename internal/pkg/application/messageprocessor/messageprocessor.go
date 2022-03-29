package messageprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/rs/zerolog"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg []byte) error
}

// hantera kö av msgs, skicka till converter registry

type msgProcessor struct {
	dmc    domain.DeviceManagementClient
	conReg conversion.ConverterRegistry
	event  events.EventPublisher
	log    zerolog.Logger
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventPublisher, log zerolog.Logger) MessageProcessor {
	return &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
		log:    log,
	}
}

func (mp *msgProcessor) ProcessMessage(ctx context.Context, msg []byte) error {
	// extract and send devEUI to device management client
	// format is from mqtt, not device management client

	dm := DeviceMessage{}

	err := json.Unmarshal(msg, &dm)
	if err == nil {
		result, err := mp.dmc.FindDeviceFromDevEUI(ctx, dm.DevEUI)
		if err == nil {
			// response with internal id, type and gets passed to Converter registry
			// converter registry returns the correct converters
			messageConverters := mp.conReg.DesignateConverters(ctx, result.Types)
			if len(messageConverters) == 0 {
				return fmt.Errorf("no matching converters for device")
			}

			for _, mc := range messageConverters {
				// msg converter converts msg payload to internal format and returns it
				payload, err := mc.ConvertPayload(ctx, result.InternalID, msg)
				if err == nil {
					mp.event.Publish(ctx, payload)
					// converted message gets sent to event publisher
				}
			}
		}
	}

	return err
}

type DeviceMessage struct {
	DevEUI string `json:"devEUI"`
}

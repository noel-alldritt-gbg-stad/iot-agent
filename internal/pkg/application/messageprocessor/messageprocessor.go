package messageprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	iotcore "github.com/diwise/iot-core/pkg/messaging/events"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg decoder.Payload) error
}

type msgProcessor struct {
	dmc    domain.DeviceManagementClient
	conReg conversion.ConverterRegistry
	event  events.EventSender
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventSender) MessageProcessor {
	return &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
	}
}

func (mp *msgProcessor) ProcessMessage(ctx context.Context, msg decoder.Payload) error {
	log := logging.GetFromContext(ctx)

	device, err := mp.dmc.FindDeviceFromDevEUI(ctx, msg.DevEUI)
	if err != nil {
		log.Error().Err(err).Msg("device lookup failure")
		return err
	}

	err = mp.event.Publish(ctx, events.NewStatusMessage(device.InternalID))
	if err != nil {
		log.Error().Err(err).Msg("failed to publish status message")
	}

	if msg.Error != "" {
		log.Info().Msg("ignoring payload due to device error")
		return nil
	}

	messageConverters := mp.conReg.DesignateConverters(ctx, device.Types)
	if len(messageConverters) == 0 {
		return fmt.Errorf("no matching converters for device")
	}

	for _, convert := range messageConverters {
		payload, err := convert(ctx, device.InternalID, msg)
		if err != nil {
			log.Error().Err(err).Msg("conversion failed")
			continue
		}

		m := iotcore.MessageReceived{
			Device:    device.InternalID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Pack:      payload,
		}

		if device.IsActive {
			err = mp.event.Send(ctx, &m)
			if err != nil {
				log.Error().Err(err).Msg("failed to send event")
			}
		}
	}

	if !device.IsActive {
		log.Warn().Str("deviceID", device.InternalID).Msg("ignoring message from inactive device")
	}

	return nil
}

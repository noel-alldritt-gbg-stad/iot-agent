package messageprocessor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	iotcore "github.com/diwise/iot-core/pkg/messaging/events"
	"github.com/rs/zerolog"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg []byte) error
}

type msgProcessor struct {
	dmc    domain.DeviceManagementClient
	conReg conversion.ConverterRegistry
	event  events.EventSender
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventSender, log zerolog.Logger) MessageProcessor {
	return &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
	}
}

func (mp *msgProcessor) ProcessMessage(ctx context.Context, msg []byte) error {
	dm := struct {
		DevEUI string `json:"devEUI"`
		Error  string `json:"error"`
		Type   string `json:"type"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return err
	}

	log := logging.GetFromContext(ctx)
	log.Info().Msgf("received payload from %s: %s", dm.DevEUI, string(msg))

	result, err := mp.dmc.FindDeviceFromDevEUI(ctx, dm.DevEUI)
	if err != nil {
		log.Error().Err(err).Msg("device lookup failure")
		return err
	}

	if dm.Error != "" {
		log.Info().Msg("ignoring payload due to device error")
		return nil
	}

	messageConverters := mp.conReg.DesignateConverters(ctx, result.Types)
	if len(messageConverters) == 0 {
		return fmt.Errorf("no matching converters for device")
	}

	for _, convert := range messageConverters {
		payload, err := convert(ctx, result.InternalID, msg)
		if err != nil {
			log.Error().Err(err).Msg("conversion failed")
			continue
		}

		m := iotcore.MessageReceived {
			Device:    result.InternalID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Pack:      payload,
		}

		err = mp.event.Send(ctx, m)
		if err != nil {
			log.Error().Err(err).Msg("failed to send event")
		}
	}

	return err
}

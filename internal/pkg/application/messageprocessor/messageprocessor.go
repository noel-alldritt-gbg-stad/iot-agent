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
	event  events.EventSender
	log    zerolog.Logger
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventSender, log zerolog.Logger) MessageProcessor {
	return &msgProcessor{
		dmc:    dmc,
		conReg: conReg,
		event:  event,
		log:    log,
	}
}

func (mp *msgProcessor) ProcessMessage(ctx context.Context, msg []byte) error {
	dm := struct {
		DevEUI string `json:"devEUI"`
		Error  string `json:"error"`
		Type   string `json:"type"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err == nil {
		mp.log.Info().Msgf("received payload from %s: %s", dm.DevEUI, string(msg))

		result, err := mp.dmc.FindDeviceFromDevEUI(ctx, dm.DevEUI)
		if err != nil {
			mp.log.Error().Err(err).Msg("device lookup failure")
			return err
		}

		if dm.Error != "" {
			mp.log.Info().Msg("ignoring payload due to device error")
			return nil
		}

		messageConverters := mp.conReg.DesignateConverters(ctx, result.Types)
		if len(messageConverters) == 0 {
			return fmt.Errorf("no matching converters for device")
		}

		for _, convert := range messageConverters {
			payload, err := convert(ctx, mp.log, result.InternalID, msg)
			if err != nil {
				mp.log.Error().Err(err).Msg("conversion failed")
				continue
			}

			justlooking, _ := json.Marshal(payload)
			mp.log.Info().Msgf("successfully converted incoming message to internal format: %s", justlooking)
			err = mp.event.Send(ctx, *payload)
			if err != nil {
				mp.log.Error().Err(err).Msg("failed to send event")
			}
		}

		return err
	}

	return err
}

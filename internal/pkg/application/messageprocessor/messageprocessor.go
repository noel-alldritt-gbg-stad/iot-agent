package messageprocessor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/logging"
	"github.com/rs/zerolog"
)

type MessageProcessor interface {
	ProcessMessage(ctx context.Context, msg []byte) error
}

// hantera kö av msgs, skicka till converter registry

type msgProcessor struct {
	dmc        domain.DeviceManagementClient
	conReg     conversion.ConverterRegistry
	event      events.EventSender
	decoderReg decoder.DecoderRegistry
}

func NewMessageReceivedProcessor(dmc domain.DeviceManagementClient, conReg conversion.ConverterRegistry, event events.EventSender, decoderReg decoder.DecoderRegistry, log zerolog.Logger) MessageProcessor {
	return &msgProcessor{
		dmc:        dmc,
		conReg:     conReg,
		event:      event,
		decoderReg: decoderReg,
	}
}

func (mp *msgProcessor) ProcessMessage(ctx context.Context, msg []byte) error {
	dm := struct {
		DevEUI     string `json:"devEUI"`
		Error      string `json:"error"`
		Type       string `json:"type"`
		SensorType string `json:"sensorType"`
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

		err = mp.event.Send(ctx, *payload)
		if err != nil {
			log.Error().Err(err).Msg("failed to send event")
		}
	}

	return err
}

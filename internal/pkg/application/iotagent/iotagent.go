package iotagent

import (
	"context"
	"encoding/json"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/iot-agent/internal/pkg/infrastructure/logging"
	"github.com/rs/zerolog"
)

type IoTAgent interface {
	MessageReceived(ctx context.Context, msg []byte) error
}

type iotAgent struct {
	mp  messageprocessor.MessageProcessor
	dr  decoder.DecoderRegistry
	dmc domain.DeviceManagementClient
}

func NewIoTAgent(dmc domain.DeviceManagementClient, eventPub events.EventSender, log zerolog.Logger) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	decreg := decoder.NewDecoderRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub, log)

	return &iotAgent{
		mp:  msgprcs,
		dr:  decreg,
		dmc: dmc,
	}
}

func (a *iotAgent) MessageReceived(ctx context.Context, msg []byte) error {
	log := logging.GetFromContext(ctx)

	devEUI, err := getDevEUIFromMessage(msg)
	if err != nil {
		log.Error().Err(err).Msg("unable to get DevEUI from payload")
		return err
	}
	
	device, err := a.dmc.FindDeviceFromDevEUI(ctx, devEUI)
	if err != nil {
		log.Error().Err(err).Msg("device lookup failure")
		return err
	}
	
	decoder := a.dr.GetDecoderForSensorType(ctx, device.SensorType)

	err = decoder(ctx, msg, func(c context.Context, m []byte) error {
		err = a.mp.ProcessMessage(c, m)
		if err != nil {
			log.Error().Err(err).Msg("failed to process message")
		}
		return err
	})

	return err
}

func getDevEUIFromMessage(msg []byte) (string, error) {
	dm := struct {
		DevEUI string `json:"devEUI"`
	}{}

	var err error

	if err = json.Unmarshal(msg, &dm); err == nil {
		return dm.DevEUI, nil
	}

	dmList := []struct {
		DevEUI string `json:"devEUI"`
	}{}

	if err = json.Unmarshal(msg, &dmList); err == nil {
		return dmList[0].DevEUI, nil
	}

	return "", err
}

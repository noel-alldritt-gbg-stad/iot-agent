package iotagent

import (
	"context"
	"encoding/json"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/rs/zerolog"
)

type IoTAgent interface {
	MessageReceived(ctx context.Context, msg []byte) error
}

type iotAgent struct {
	mp  messageprocessor.MessageProcessor
	dr  decoder.DecoderRegistry
	log zerolog.Logger
}

func NewIoTAgent(dmc domain.DeviceManagementClient, eventPub events.EventSender, log zerolog.Logger) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	decreg := decoder.NewDecoderRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub, log)

	return &iotAgent{
		mp:  msgprcs,
		dr:  decreg,
		log: log,
	}
}

func (a *iotAgent) MessageReceived(ctx context.Context, msg []byte) error {

	sensorType, err := parseSensorType(msg)
	if err != nil {
		return err
	}

	dfn := a.dr.GetDecodersForSensorType(ctx, sensorType)

	err = dfn(ctx, msg, func(c context.Context, m []byte) error {
		err = a.mp.ProcessMessage(c, m)
		if err != nil {
			a.log.Error().Err(err).Msg("failed to process message")
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func parseSensorType(msg []byte) (string, error) {
	dm := struct {
		SensorType        string `json:"sensorType"`
		DeviceProfileName string `json:"deviceProfileName"`
	}{}

	dmA := []struct {
		SensorType        string `json:"sensorType"`
		DeviceProfileName string `json:"deviceProfileName"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		err = json.Unmarshal(msg, &dmA)
		if err != nil {
			return "", err
		}
		if dmA[0].SensorType != "" {
			return dmA[0].SensorType, nil
		}

		if dmA[0].DeviceProfileName != "" {
			return dmA[0].DeviceProfileName, nil
		}
	}

	if dm.SensorType != "" {
		return dm.SensorType, nil
	}

	if dm.DeviceProfileName != "" {
		return dm.DeviceProfileName, nil
	}

	return "", nil
}

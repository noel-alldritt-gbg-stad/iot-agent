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
	mp messageprocessor.MessageProcessor
	dr decoder.DecoderRegistry
	log zerolog.Logger
}

func NewIoTAgent(dmc domain.DeviceManagementClient, eventPub events.EventSender, log zerolog.Logger) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	decreg := decoder.NewDecoderRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub, decreg, log)

	return &iotAgent{
		mp: msgprcs,
		dr: decreg,
		log: log,
	}
}

func (a *iotAgent) MessageReceived(ctx context.Context, msg []byte) error {

	dm := struct {
		SensorType string `json:"sensorType"`
	}{}
	
	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return err
	}

	dfn := a.dr.GetDecodersForSensorType(ctx, dm.SensorType)	

	err = dfn(ctx, msg, func (context.Context, []byte) error {
		err = a.mp.ProcessMessage(ctx, msg)
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

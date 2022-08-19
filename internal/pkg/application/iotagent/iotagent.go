package iotagent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/application/messageprocessor"
	dmc "github.com/diwise/iot-device-mgmt/pkg/client"
)

//go:generate moq -rm -out iotagent_mock.go . IoTAgent

type IoTAgent interface {
	MessageReceived(ctx context.Context, msg []byte) error
}

type iotAgent struct {
	mp  messageprocessor.MessageProcessor
	dr  decoder.DecoderRegistry
	dmc dmc.DeviceManagementClient
}

func NewIoTAgent(dmc dmc.DeviceManagementClient, eventPub events.EventSender) IoTAgent {
	conreg := conversion.NewConverterRegistry()
	decreg := decoder.NewDecoderRegistry()
	msgprcs := messageprocessor.NewMessageReceivedProcessor(dmc, conreg, eventPub)

	return &iotAgent{
		mp:  msgprcs,
		dr:  decreg,
		dmc: dmc,
	}
}

func (a *iotAgent) MessageReceived(ctx context.Context, msg []byte) error {

	devEUI, err := getDevEUIFromMessage(msg)
	if err != nil {
		return fmt.Errorf("unable to get DevEUI from payload (%w)", err)
	}

	device, err := a.dmc.FindDeviceFromDevEUI(ctx, devEUI)
	if err != nil {
		return fmt.Errorf("device lookup failure (%w)", err)
	}

	d := a.dr.GetDecoderForSensorType(ctx, device.SensorType())

	err = d(ctx, msg, func(c context.Context, m decoder.Payload) error {
		err = a.mp.ProcessMessage(c, m)
		if err != nil {
			err = fmt.Errorf("failed to process message (%w)", err)
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

package iotagent

import (
	"context"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/farshidtz/senml/v2"
	"github.com/farshidtz/senml/v2/codec"
	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestSenlabTPayload(t *testing.T) {
	is, dmc, e, log := testSetup(t)

	app := NewIoTAgent(dmc, e, log)
	err := app.MessageReceived(context.Background(), []byte(senlabT))

	is.NoErr(err)
	is.True(len(e.SendCalls()) > 0)

	pack, err := codec.Decode(senml.MediaTypeSenmlJSON, e.SendCalls()[0].Msg)
	is.NoErr(err)
	is.True(*pack[1].Value == 6.625)
}

func TestStripsPayload(t *testing.T) {
	is, dmc, e, log := testSetup(t)

	app := NewIoTAgent(dmc, e, log)
	err := app.MessageReceived(context.Background(), []byte(stripsPayload))

	is.NoErr(err)
	is.True(len(e.SendCalls()) > 0)

	pack, err := codec.Decode(senml.MediaTypeSenmlJSON, e.SendCalls()[0].Msg)
	is.NoErr(err)
	is.True(pack[0].BaseName == "urn:oma:lwm2m:ext:3303")
}

func TestElsysPayload(t *testing.T) {
	is, dmc, e, log := testSetup(t)

	app := NewIoTAgent(dmc, e, log)
	err := app.MessageReceived(context.Background(), []byte(elsys))

	is.NoErr(err)
	is.True(len(e.SendCalls()) > 0)

	pack, err := codec.Decode(senml.MediaTypeSenmlJSON, e.SendCalls()[0].Msg)
	is.NoErr(err)
	is.True(*pack[1].Value == 19.3)
}

func TestErsPayload(t *testing.T) {
	is, dmc, e, log := testSetup(t)

	app := NewIoTAgent(dmc, e, log)
	err := app.MessageReceived(context.Background(), []byte(ers))

	is.NoErr(err)
	is.True(len(e.SendCalls()) > 0)

	pack, err := codec.Decode(senml.MediaTypeSenmlJSON, e.SendCalls()[0].Msg)
	is.NoErr(err)
	is.True(*pack[1].Value == 23.8)
}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, *events.EventSenderMock, zerolog.Logger) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (*domain.Result, error) {

			res := &domain.Result{
				InternalID: "internal-id-for-device",
				Types:      []string{"urn:oma:lwm2m:ext:3303"},
			}

			if devEUI == "70b3d580a010f260" {
				res.SensorType = "tem_lab_14ns"
			} else if devEUI == "70b3d52c00019193" {
				res.SensorType = "strips_lora_ms_h"
			} else {
				res.SensorType = "Elsys_Codec"
			}

			return res, nil
		},
	}

	e := &events.EventSenderMock{
		SendFunc: func(ctx context.Context, msg []byte) error {
			return nil
		},
	}

	return is, dmc, e, zerolog.Logger{}
}

const senlabT string = `[{
	"devEui": "70b3d580a010f260",
	"sensorType": "tem_lab_14ns",
	"timestamp": "2022-04-12T05:08:50.301732Z",
	"payload": "01FE90619c10006A",
	"spreadingFactor": 12,
	"rssi": "-113",
	"snr": "-11.8",
	"gatewayIdentifier": 184,
	"fPort": "3",
	"latitude": 57.806266,
	"longitude": 12.07727
}]`

const stripsPayload string = `[{"devEui":"70b3d52c00019193","sensorType":"strips_lora_ms_h","timestamp":"2022-04-21T09:33:40.713643Z","payload":"ffff01590200d90400d4063c07000008000009000a01","spreadingFactor":"10","rssi":"-108","snr":"-3","gatewayIdentifier":"824","fPort":"1"}]`

const elsys string = `{
	"applicationID": "8",
	"applicationName": "Water-Temperature",
	"deviceName": "sk-elt-temp-16",
	"deviceProfileName": "Elsys_Codec",
	"deviceProfileID": "xxxxxxxxxxxx",
	"devEUI": "xxxxxxxxxxxxxx",
	"rxInfo": [{
		"gatewayID": "xxxxxxxxxxx",
		"uplinkID": "xxxxxxxxxxx",
		"name": "SN-LGW-047",
		"time": "2022-03-28T12:40:40.653515637Z",
		"rssi": -105,
		"loRaSNR": 8.5,
		"location": {
			"latitude": 62.36956091265246,
			"longitude": 17.319844410529534,
			"altitude": 0
		}
	}],
	"txInfo": {
		"frequency": 867700000,
		"dr": 5
	},
	"adr": true,
	"fCnt": 10301,
	"fPort": 5,
	"data": "Bw2KDADB",
	"object": {
		"externalTemperature": 19.3,
		"vdd": 3466
	},
	"tags": {
		"Location": "Vangen"
	}
}`

const ers string = `
{
    "deviceName": "mcg-ers-co2-01",
    "deviceProfileName": "ELSYS",
    "deviceProfileID": "0b765672-274a-41eb-b1c5-bb2bec9d14e8",
    "devEUI": "a81758fffe05e6fb",
    "data": "AQDuAhYEALIFAgYBxAcONA==",
    "object": {
        "co2": 452,
        "humidity": 22,
        "light": 178,
        "motion": 2,
        "temperature": 23.8,
        "vdd": 3636
    }
}`

package decoder

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestSenlabTBasicDecoder(t *testing.T) {
	is, _ := testSetup(t)

	r := &Payload{}

	err := SenlabTBasicDecoder(context.Background(), []byte(senlabT), func(c context.Context, m []byte) error {
		json.Unmarshal(m, &r)
		return nil
	})

	is.True(r.Timestamp == "2022-04-12T05:08:50.301732Z")
	is.NoErr(err)
}

func TestElsysDecoder(t *testing.T) {
	is, _ := testSetup(t)

	r := &Payload{}

	err := ElsysDecoder(context.Background(), []byte(elsys), func(c context.Context, m []byte) error {
		json.Unmarshal(m, &r)
		return nil
	})

	is.True(r.SensorType == "Elsys_Codec")
	is.NoErr(err)
}

func TestSenlabTBasicDecoderSensorReadingError(t *testing.T) {
	is, _ := testSetup(t)

	err := SenlabTBasicDecoder(context.Background(), []byte(senlabT_sensorReadingError), func(c context.Context, m []byte) error {
		return nil
	})

	is.True(err != nil)
}

func testSetup(t *testing.T) (*is.I, zerolog.Logger) {
	is := is.New(t)
	return is, zerolog.Logger{}
}

const senlabT string = `[{
	"devEui": "70b3d580a010f260",
	"sensorType": "tem_lab_14ns",
	"timestamp": "2022-04-12T05:08:50.301732Z",
	"payload": "01FE90619c10006A",
	"spreadingFactor": 12,
	"rssi": -113,
	"snr": -11.8,
	"gatewayIdentifier": 184,
	"fPort": 3,
	"latitude": 57.806266,
	"longitude": 12.07727
}]`

// payload ...0xFD14 = -46.75 = sensor reading error
const senlabT_sensorReadingError string = `[{
	"devEui": "70b3d580a010f260",
	"sensorType": "tem_lab_14ns",
	"timestamp": "2022-04-12T05:08:50.301732Z",
	"payload": "01FE90619c10FD14",
	"spreadingFactor": 12,
	"rssi": -113,
	"snr": -11.8,
	"gatewayIdentifier": 184,
	"fPort": 3,
	"latitude": 57.806266,
	"longitude": 12.07727
}]`

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
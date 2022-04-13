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

	r := struct {
		DevEUI       string
		Id           int
		BatteryLevel int
		Temperature  float32
		Timestamp    string
	}{}

	err := SenlabTBasicDecoder(context.Background(), []byte(senlabT), func(c context.Context, m []byte) error {
		json.Unmarshal(m, &r)
		return nil
	})

	is.True(r.Id == 1)	
	is.True(r.BatteryLevel == 100)
	is.True(r.Temperature == 6.625)
	is.True(r.Timestamp == "2022-04-12T05:08:50.301732Z")
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

package conversion

import (
	"context"
	"testing"
	"time"

	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/matryer/is"
)

func TestThatTemperatureDecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := decoder.Payload{
		DevEUI:    "ncaknlclkdanklcd",
		Timestamp: "2006-01-02T15:04:05Z",
	}
	temp := struct {
		Temperature float32 `json:"temperature"`
	}{
		22.2,
	}
	payload.Measurements = append(payload.Measurements, temp)

	msg, err := Temperature(ctx, "internalID", payload)

	is.NoErr(err)
	is.Equal(22.2, *msg[1].Value)
}

func TestThatCO2DecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := decoder.Payload{
		DevEUI:    "ncaknlclkdanklcd",
		Timestamp: "2006-01-02T15:04:05Z",
	}
	co2 := struct {
		CO2 int `json:"co2"`
	}{
		22,
	}
	payload.Measurements = append(payload.Measurements, co2)

	msg, err := AirQuality(ctx, "internalID", payload)

	is.NoErr(err)
	is.Equal(22.0, *msg[1].Value)
}

func TestThatPresenceDecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := decoder.Payload{
		DevEUI:    "ncaknlclkdanklcd",
		Timestamp: "2006-01-02T15:04:05Z",
	}
	present := struct {
		Presence bool `json:"present"`
	}{
		true,
	}
	payload.Measurements = append(payload.Measurements, present)

	msg, err := Presence(ctx, "internalID", payload)

	is.NoErr(err)
	is.True(*msg[1].BoolValue)
}

func TestThatWatermeterDecodesValuesCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)

	payload := decoder.Payload{
		DevEUI:     "3489573498573459",
		DeviceName: "deviceName",
		Timestamp:  time.Now().Format(time.RFC3339),
	}
	curDateTime := struct {
		CurrentDateTime string `json:"curDateTime"`
	}{
		"2006-01-02T15:04:05Z",
	}
	payload.Measurements = append(payload.Measurements, curDateTime)
	curVol := struct {
		CurrentVolume float64 `json:"curVol"`
	}{
		1009,
	}
	payload.Measurements = append(payload.Measurements, curVol)

	msg, err := Watermeter(ctx, "internalID", payload)

	is.NoErr(err)
	is.True(msg != nil)
}

func TestThatWatermeterCanBeDecodedAndConverted(t *testing.T){
	is, ctx := mcmTestSetup(t)

	err := decoder.WatermeteringDecoder(ctx, []byte(watermetering), watermeteringMessageProcessor)

	is.NoErr(err)
}

func watermeteringMessageProcessor(ctx context.Context, msg decoder.Payload) error {
	_, err := Watermeter(ctx, "", msg)
	return err
}

func mcmTestSetup(t *testing.T) (*is.I, context.Context) {
	ctx, _ := logging.NewLogger(context.Background(), "test", "")
	return is.New(t), ctx
}

const watermetering string = `
{
  "applicationID": "2",
  "applicationName": "Watermetering",
  "deviceName": "05394167",
  "deviceProfileName": "Axioma_Universal_Codec",
  "deviceProfileID": "8be301da",
  "devEUI": "3489573498573459",
  "txInfo": { "frequency": 867100000, "dr": 0 },
  "adr": true,
  "fCnt": 182,
  "fPort": 100,
  "data": "//8VAQ==",
  "object": {
    "curDateTime": "2022-02-10 15:13:57",
    "curVol": 1009,
    "deltaVol": {
      "id1": 0,
      "id10": 13,
      "id11": 10,
      "id12": 2,
      "id13": 0,
      "id14": 1,
      "id15": 0,
      "id16": 5,
      "id17": 0,
      "id18": 0,
      "id19": 0,
      "id2": 8,
      "id20": 2,
      "id21": 0,
      "id22": 0,
      "id23": 0,
      "id3": 0,
      "id4": 0,
      "id5": 0,
      "id6": 0,
      "id7": 0,
      "id8": 5,
      "id9": 6
    },
    "frameVersion": 1,
    "statusCode": 0
  },
  "tags": { "Location": "UnSet", "SerialNo": "05394167" }
}
`
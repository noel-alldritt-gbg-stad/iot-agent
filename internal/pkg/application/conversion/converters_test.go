package conversion

import (
	"context"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/matryer/is"
)

func TestThatTemperatureDecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := decoder.Payload{
		DevEUI: "ncaknlclkdanklcd",
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
		DevEUI: "ncaknlclkdanklcd",
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

func mcmTestSetup(t *testing.T) (*is.I, context.Context) {
	ctx, _ := logging.NewLogger(context.Background(), "test", "")
	return is.New(t), ctx
}

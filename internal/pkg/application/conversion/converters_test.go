package conversion

import (
	"context"
	"testing"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/matryer/is"
)

func TestThatTemperatureDecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := `{"devEUI":"ncaknlclkdanklcd","timestamp":"2006-01-02T15:04:05Z","measurements":[{"temperature":22.2}]}`

	msg, err := Temperature(ctx, "internalID", []byte(payload))

	is.NoErr(err)
	is.Equal(22.2, *msg[1].Value)
}

func TestThatCO2DecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := `{"devEUI":"ncaknlclkdanklcd","timestamp":"2006-01-02T15:04:05Z","measurements":[{"co2":22}]}`

	msg, err := AirQuality(ctx, "internalID", []byte(payload))

	is.NoErr(err)
	is.Equal(22.0, *msg[1].Value)
}

func mcmTestSetup(t *testing.T) (*is.I, context.Context) {
	ctx, _ := logging.NewLogger(context.Background(), "test", "")
	return is.New(t), ctx
}

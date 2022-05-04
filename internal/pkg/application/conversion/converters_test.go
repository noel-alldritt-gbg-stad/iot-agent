package conversion

import (
	"context"
	"testing"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/matryer/is"
)

func TestThatTemperatureDecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := `{"devEUI":"ncaknlclkdanklcd","measurements":[{"temperature":22.2}]}`

	msg, err := Temperature(ctx, "internalID", []byte(payload))

	is.NoErr(err)
	is.Equal(`[{"bn":"urn:oma:lwm2m:ext:3303","n":"0","vs":"internalID"},{"n":"Temperature","v":22.2}]`, string(msg))
}

func TestThatCO2DecodesValueCorrectly(t *testing.T) {
	is, ctx := mcmTestSetup(t)
	payload := `{"devEUI":"ncaknlclkdanklcd","measurements":[{"co2":22}]}`

	msg, err := AirQuality(ctx, "internalID", []byte(payload))

	is.NoErr(err)
	is.Equal(`[{"bn":"urn:oma:lwm2m:ext:3428","n":"0","vs":"internalID"},{"n":"CO2","v":22}]`, string(msg))
}

func mcmTestSetup(t *testing.T) (*is.I, context.Context) {
	ctx, _ := logging.NewLogger(context.Background(), "test", "")
	return is.New(t), ctx
}

package conversion

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestThatConvertPayloadGetsTemperatureValueCorrectly(t *testing.T) {
	is, log := mcmTestSetup(t)

	mc := &msgConverter{
		Type: "urn:oma:lwm2m:ext:3303",
	}

	payload := `{"devEUI":"ncaknlclkdanklcd","object":{"externalTemperature":22.2}}`

	msg, err := mc.ConvertPayload(context.Background(), log, "internalID", []byte(payload))
	is.NoErr(err)
	is.Equal(msg.SensorValue, 22.2)
}

func mcmTestSetup(t *testing.T) (*is.I, zerolog.Logger) {
	return is.New(t), zerolog.Logger{}
}

package conversion

import (
	"context"
	"testing"

	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestThatTemperatureDecodesValueCorrectly(t *testing.T) {
	is, log := mcmTestSetup(t)
	payload := `{"devEUI":"ncaknlclkdanklcd","object":{"externalTemperature":22.2}}`

	msg, err := Temperature(context.Background(), log, "internalID", []byte(payload))

	is.NoErr(err)
	is.Equal(msg.SensorValue, 22.2)
}

func mcmTestSetup(t *testing.T) (*is.I, zerolog.Logger) {
	return is.New(t), zerolog.Logger{}
}

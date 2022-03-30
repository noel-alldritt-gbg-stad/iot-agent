package messageprocessor

import (
	"context"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestFailsOnInvalidMessage(t *testing.T) {
	is, dmc, cr, ep, log := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, ep, log)

	err := mp.ProcessMessage(context.Background(), []byte("msg"))
	is.True(err != nil)
}

func TestProcessMessageWorksWithCorrectInput(t *testing.T) {
	is, dmc, cr, ep, log := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, ep, log)

	err := mp.ProcessMessage(context.Background(), []byte(payload))
	is.NoErr(err)
}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, conversion.ConverterRegistry, events.EventPublisher, zerolog.Logger) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (domain.Result, error) {
			return domain.Result{
					InternalID: "internalID",
					Types:      []string{"urn:oma:lwm2m:ext:3303"},
				},
				nil
		},
	}
	cr := &conversion.ConverterRegistryMock{
		DesignateConvertersFunc: func(ctx context.Context, types []string) []conversion.MessageConverter {
			return []conversion.MessageConverter{
				&conversion.MessageConverterMock{
					ConvertPayloadFunc: func(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (conversion.InternalMessageFormat, error) {
						return conversion.InternalMessageFormat{}, nil
					},
				},
			}
		},
	}
	ep := &events.EventPublisherMock{
		PublishFunc: func(ctx context.Context, msg conversion.InternalMessageFormat) error {
			return nil
		},
	}

	return is, dmc, cr, ep, zerolog.Logger{}
}

const payload string = `{"devEUI":"xxxxxxxxxxxxxx","object":{"externalTemperature":23.5}}`

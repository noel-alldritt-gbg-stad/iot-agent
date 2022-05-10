package messageprocessor

import (
	"context"
	"errors"
	"testing"	

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	iotcore "github.com/diwise/iot-core/pkg/messaging/events"
	"github.com/farshidtz/senml/v2"
	"github.com/matryer/is"
	"github.com/rs/zerolog"
)

func TestFailsOnInvalidMessage(t *testing.T) {
	is, dmc, cr, ep, log := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, ep, log)

	err := mp.ProcessMessage(context.Background(), []byte("msg"))
	is.True(err != nil)
}

func TestFailsOnInvalidType(t *testing.T) {
	is, _, cr, ep, log := testSetup(t)

	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (*domain.Result, error) {
			return &domain.Result{}, errors.New("devEUI does not belong to a sensor of any valid types")
		},
	}

	mp := NewMessageReceivedProcessor(dmc, cr, ep, log)

	err := mp.ProcessMessage(context.Background(), []byte(payload))
	is.True(err != nil)
	is.Equal(err.Error(), "devEUI does not belong to a sensor of any valid types")
}

func TestProcessMessageWorksWithValidTemperatureInput(t *testing.T) {
	is, dmc, cr, ep, log := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, ep, log)

	err := mp.ProcessMessage(context.Background(), []byte(payload))
	is.NoErr(err)
}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, conversion.ConverterRegistry, events.EventSender, zerolog.Logger) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (*domain.Result, error) {
			return &domain.Result{
				InternalID: "internalID",
				Types:      []string{"urn:oma:lwm2m:ext:3303"},
			}, nil
		},
	}
	cr := &conversion.ConverterRegistryMock{
		DesignateConvertersFunc: func(ctx context.Context, types []string) []conversion.MessageConverterFunc {
			return []conversion.MessageConverterFunc{
				func(ctx context.Context, internalID string, msg []byte) (senml.Pack, error) {
					return senml.Pack{}, nil
				},
			}
		},
	}
	ep := &events.EventSenderMock{
		SendFunc: func(ctx context.Context, m iotcore.MessageReceived) error {
			return nil
		},
	}

	return is, dmc, cr, ep, zerolog.Logger{}
}

const payload string = `{"devEUI":"xxxxxxxxxxxxxx","object":{"externalTemperature":23.5}}`

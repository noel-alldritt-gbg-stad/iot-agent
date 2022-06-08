package messageprocessor

import (
	"context"
	"errors"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/diwise/messaging-golang/pkg/messaging"
	"github.com/farshidtz/senml/v2"
	"github.com/matryer/is"
)

func TestFailsOnInvalidType(t *testing.T) {
	is, _, cr, ep := testSetup(t)

	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (*domain.Result, error) {
			return &domain.Result{}, errors.New("devEUI does not belong to a sensor of any valid types")
		},
	}

	mp := NewMessageReceivedProcessor(dmc, cr, ep)

	err := mp.ProcessMessage(context.Background(), newPayload())
	is.True(err != nil)
	is.Equal(err.Error(), "devEUI does not belong to a sensor of any valid types")
}

func TestProcessMessageWorksWithValidTemperatureInput(t *testing.T) {
	is, dmc, cr, ep := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, ep)

	err := mp.ProcessMessage(context.Background(), newPayload())
	is.NoErr(err)
	is.Equal(len(ep.SendCalls()), 1) // should have been called once
}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, conversion.ConverterRegistry, *events.EventSenderMock) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (*domain.Result, error) {
			return &domain.Result{
				InternalID: "internalID",
				Types:      []string{"urn:oma:lwm2m:ext:3303"},
				IsActive:   true,
			}, nil
		},
	}
	cr := &conversion.ConverterRegistryMock{
		DesignateConvertersFunc: func(ctx context.Context, types []string) []conversion.MessageConverterFunc {
			return []conversion.MessageConverterFunc{
				func(ctx context.Context, internalID string, payload decoder.Payload) (senml.Pack, error) {
					return senml.Pack{}, nil
				},
			}
		},
	}
	ep := &events.EventSenderMock{
		SendFunc: func(ctx context.Context, m messaging.CommandMessage) error {
			return nil
		},
		PublishFunc: func(ctx context.Context, m messaging.TopicMessage) error {
			return nil
		},
	}

	return is, dmc, cr, ep
}

func newPayload() decoder.Payload {
	payload := decoder.Payload{
		DevEUI:    "ncaknlclkdanklcd",
		Timestamp: "2006-01-02T15:04:05Z",
	}
	temp := struct {
		Temperature float32 `json:"temperature"`
	}{
		23.5,
	}
	payload.Measurements = append(payload.Measurements, temp)

	return payload
}

package messageprocessor

import (
	"context"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/matryer/is"
)

func TestFailsOnInvalidMessage(t *testing.T) {
	is, dmc, _, _ := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, nil, nil)

	err := mp.ProcessMessage(context.Background(), []byte("msg"))
	is.True(err != nil)
}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, *conversion.ConverterRegistry, *events.EventPublisher) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (domain.Result, error) {
			return domain.Result{},
				nil
		},
	}

	return is, dmc, nil, nil
}

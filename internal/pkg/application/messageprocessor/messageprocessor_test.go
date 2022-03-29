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
	is, dmc, cr, _ := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, nil)

	err := mp.ProcessMessage(context.Background(), []byte("msg"))
	is.True(err != nil)
}

func TestXxx(t *testing.T) {
	is, dmc, cr, _ := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, cr, nil)

	err := mp.ProcessMessage(context.Background(), []byte(payload))
	is.True(err != nil)

}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, conversion.ConverterRegistry, *events.EventPublisher) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(ctx context.Context, devEUI string) (domain.Result, error) {
			return domain.Result{
					InternalID: "internalID",
					Types:      []string{"temperature"},
				},
				nil
		},
	}
	cr := &conversion.ConverterRegistryMock{
		DesignateConvertersFunc: func(ctx context.Context, types []string) []conversion.MessageConverter {
			return []conversion.MessageConverter{}
		},
	}

	return is, dmc, cr, nil
}

const payload string = `{"applicationID":"8","applicationName":"Water-Temperature","deviceName":"sk-elt-temp-16","deviceProfileName":"Elsys_Codec","deviceProfileID":"xxxxxxxxxxxx","devEUI":"xxxxxxxxxxxxxx"}`

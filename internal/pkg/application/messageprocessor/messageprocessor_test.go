package messageprocessor

import (
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/matryer/is"
)

func TestXxx(t *testing.T) {

}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(devEUI string) (domain.Result, error) {
			return domain.Result{
					InternalID: "internalID",
					Types: []string{
						"temperature", "water",
					},
				},
				nil
		},
	}

	return is, dmc
}

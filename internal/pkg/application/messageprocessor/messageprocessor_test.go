package messageprocessor

import (
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application/conversion"
	"github.com/diwise/iot-agent/internal/pkg/application/events"
	"github.com/diwise/iot-agent/internal/pkg/domain"
	"github.com/matryer/is"
)

func TestFailsOnInvalidMessage(t *testing.T) {
	is, dmc, _, _ := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, nil, nil)

	err := mp.ProcessMessage([]byte("msg"))
	is.True(err != nil)
}

func TestXxx(t *testing.T) {
	is, dmc, _, _ := testSetup(t)
	mp := NewMessageReceivedProcessor(dmc, nil, nil)

	err := mp.ProcessMessage([]byte(payload))
	is.True(err != nil)

}

func testSetup(t *testing.T) (*is.I, *domain.DeviceManagementClientMock, *conversion.ConverterRegistry, *events.EventPublisher) {
	is := is.New(t)
	dmc := &domain.DeviceManagementClientMock{
		FindDeviceFromDevEUIFunc: func(devEUI string) (domain.Result, error) {
			return domain.Result{
					InternalID: "internalID",
					Types:      []string{"watertemp"},
				},
				nil
		},
	}

	return is, dmc, nil, nil
}

const payload string = `{"level":"info","service":"iot-agent","version":"","mqtt-host":"iot.serva.net","timestamp":"2022-03-28T14:39:11.695538+02:00","message":"received payload: {\"applicationID\":\"8\",\"applicationName\":\"Water-Temperature\",\"deviceName\":\"sk-elt-temp-16\",\"deviceProfileName\":\"Elsys_Codec\",\"deviceProfileID\":\"xxxxxxxxxxxx\",\"devEUI\":\"xxxxxxxxxxxxxx\",\"rxInfo\":[{\"gatewayID\":\"xxxxxxxxxxx\",\"uplinkID\":\"xxxxxxxxxxx\",\"name\":\"SN-LGW-047\",\"time\":\"2022-03-28T12:40:40.653515637Z\",\"rssi\":-105,\"loRaSNR\":8.5,\"location\":{\"latitude\":62.36956091265246,\"longitude\":17.319844410529534,\"altitude\":0}}],\"txInfo\":{\"frequency\":867700000,\"dr\":5},\"adr\":true,\"fCnt\":10301,\"fPort\":5,\"data\":\"Bw2KDADB\",\"object\":{\"externalTemperature\":19.3,\"vdd\":3466},\"tags\":{\"Location\":\"Vangen\"}}"}`

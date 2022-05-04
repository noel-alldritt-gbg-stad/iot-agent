package conversion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/farshidtz/senml/v2"
	"github.com/farshidtz/senml/v2/codec"
)

type MessageConverterFunc func(ctx context.Context, internalID string, msg []byte) (*InternalMessage, error)

func Temperature(ctx context.Context, deviceID string, msg []byte) (*InternalMessage, error) {
	dm := struct {
		Measurements []struct {
			Temp *float64 `json:"temperature"`
		} `json:"measurements"`
	}{}

	if err := json.Unmarshal(msg, &dm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal measurements: %s", err.Error())
	}

	for _, m := range dm.Measurements {
		if m.Temp != nil {
			payload := &InternalMessage{
				InternalID:  deviceID,
				Type:        "urn:oma:lwm2m:ext:3303",
				SensorValue: *m.Temp,
			}

			//TODO: range and call func?
			return payload, nil
		}
	}

	return nil, fmt.Errorf("no temperature value found in payload")
}

func AirQuality(ctx context.Context, deviceID string, msg []byte) ([]byte, error) {
	dm := struct {
		Measurements []struct {
			CO2 *int `json:"co2"`
		} `json:"measurements"`
	}{}

	if err := json.Unmarshal(msg, &dm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal measurements: %s", err.Error())
	}

	var pack senml.Pack

	pack = append(pack, senml.Record{
		BaseName:    "urn:oma:lwm2m:ext:3428",
		Name:        "0",
		StringValue: deviceID,
	})

	for _, m := range dm.Measurements {
		if m.CO2 != nil {
			co2 := float64(*m.CO2)
			rec := senml.Record{
				Name:  "CO2",
				Value: &co2,
			}

			pack = append(pack, rec)
		}
	}

	data, err := codec.EncodeJSON(pack)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal measurements: %s", err.Error())
	}

	//each measurement from payload is a new record on pack

	return data, nil
}

type InternalMessage struct {
	InternalID   string  `json:"internalID"`
	Type         string  `json:"type"`
	SensorValue  float64 `json:"sensorValue"`
	ResourceName string  `json:"resourceName"`
}

func (im InternalMessage) ContentType() string {
	return "application/json" // TODO: Decide a proper content type here
}

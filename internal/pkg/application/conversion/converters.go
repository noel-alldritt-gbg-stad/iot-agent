package conversion

import (
	"context"
	"encoding/json"
	"fmt"
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

type InternalMessage struct {
	InternalID  string  `json:"internalID"`
	Type        string  `json:"type"`
	SensorValue float64 `json:"sensorValue"`
}

func (im InternalMessage) ContentType() string {
	return "application/json" // TODO: Decide a proper content type here
}

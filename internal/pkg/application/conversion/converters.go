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

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal temperature measurements: %s", err.Error())
	}

	if len(dm.Measurements) == 0 || dm.Measurements[0].Temp == nil {
		return nil, fmt.Errorf("no temperature value found in payload")
	}

	//TODO: range and call func?

	payload := &InternalMessage{
		InternalID:  deviceID,
		Type:        "urn:oma:lwm2m:ext:3303",
		SensorValue: *dm.Measurements[0].Temp,
	}

	return payload, nil
}

type InternalMessage struct {
	InternalID  string  `json:"internalID"`
	Type        string  `json:"type"`
	SensorValue float64 `json:"sensorValue"`
}

func (im InternalMessage) ContentType() string {
	return "application/json" // TODO: Decide a proper content type here
}

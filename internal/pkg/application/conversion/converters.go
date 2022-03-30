package conversion

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

type MessageConverterFunc func(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (*InternalMessage, error)

// konvertera payload till internt format.

func Temperature(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (*InternalMessage, error) {
	dm := struct {
		Object struct {
			Temperature float64 `json:"externalTemperature"`
		} `json:"object"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return nil, err
	}

	payload := &InternalMessage{
		InternalID:  internalID,
		Type:        "urn:oma:lwm2m:ext:3303",
		SensorValue: dm.Object.Temperature,
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

func (im InternalMessage) TopicName() string {
	return "temperature"
}

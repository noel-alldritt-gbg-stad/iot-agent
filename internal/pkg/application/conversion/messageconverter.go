package conversion

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
)

type MessageConverter interface {
	ConvertPayload(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (*InternalMessage, error)
}

// konvertera payload till internt format.

type msgConverter struct {
	Type string // determines what type of data we're converting, i.e. water or air temperature etc.
}

func (mc *msgConverter) ConvertPayload(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (*InternalMessage, error) {
	if mc.Type == "urn:oma:lwm2m:ext:3303" {
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
			Type:        mc.Type,
			SensorValue: dm.Object.Temperature,
		}
		return payload, nil
	}

	return nil, fmt.Errorf("failed to convert payload, type %s is unknown", mc.Type)
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

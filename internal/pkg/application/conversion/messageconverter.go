package conversion

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

type MessageConverter interface {
	ConvertPayload(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (InternalMessageFormat, error)
}

// konvertera payload till internt format.

type msgConverter struct {
	Type string // determines what type of data we're converting, i.e. water or air temperature etc.
}

func (mc *msgConverter) ConvertPayload(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (InternalMessageFormat, error) {
	dm := &DeviceMessage{}
	err := json.Unmarshal(msg, dm)
	if err == nil {
		if mc.Type == "urn:oma:lwm2m:ext:3303" {
			payload := InternalMessageFormat{
				InternalID: internalID,
				Type:       mc.Type,
				Value:      dm.Object.ExternalTemperature,
			}
			return payload, nil
		}
	}

	return InternalMessageFormat{}, err
}

type InternalMessageFormat struct {
	InternalID string  `json:"internalID"`
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
}

type DeviceMessage struct {
	DevEUI string `json:"devEUI"`
	Object Object `json:"object"`
}

type Object struct {
	ExternalTemperature float64 `json:"externalTemperature"`
}

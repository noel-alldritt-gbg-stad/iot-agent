package conversion

import (
	"context"
)

type MessageConverter interface {
	ConvertPayload(ctx context.Context, internalID string, msg []byte) (InternalMessageFormat, error)
}

//konvertera payload till internt format

type msgConverter struct {
	Type string //determines what type of data we're converting, i.e. water or air temperature etc.
}

func (mc *msgConverter) ConvertPayload(ctx context.Context, internalID string, msg []byte) (InternalMessageFormat, error) {
	payload := InternalMessageFormat{
		InternalID: internalID,
		Type:       mc.Type,
	}

	var value float64

	if mc.Type == "urn:oma:lwm2m:ext:3303" {
		payload.Value = Value{
			Temperature: value,
		}
	}

	return payload, nil
}

type InternalMessageFormat struct {
	InternalID string
	Type       string
	Value      Value
}

type Value struct {
	Temperature float64
}

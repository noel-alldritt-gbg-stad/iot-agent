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
	imf := InternalMessageFormat{
		InternalID: internalID,
		Type:       mc.Type,
	}

	return imf, nil
}

type InternalMessageFormat struct {
	InternalID string
	Type       string
	Value      string
}

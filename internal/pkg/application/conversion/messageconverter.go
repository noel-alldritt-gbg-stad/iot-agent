package conversion

import "context"

type MessageConverter interface {
	ConvertPayload(ctx context.Context, msg []byte) (InternalMessageFormat, error)
}

//konvertera payload till internt format

type msgConverter struct {
	Type string
}

func (mc *msgConverter) ConvertPayload(ctx context.Context, msg []byte) (InternalMessageFormat, error) {
	return InternalMessageFormat{}, nil
}

type InternalMessageFormat struct {
}

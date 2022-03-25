package conversion

type MessageConverter interface {
	ConvertPayload(msg []byte) (InternalMessageFormat, error)
}

//konvertera payload till internt format

type msgConverter struct {
}

func (mc *msgConverter) ConvertPayload(msg []byte) (InternalMessageFormat, error) {
	return InternalMessageFormat{}, nil
}

type InternalMessageFormat struct {
}

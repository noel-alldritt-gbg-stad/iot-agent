package decoder

import (
	"context"
)

type DecoderRegistry interface {
	DesignateDecoders(ctx context.Context, sensorType string) MessageDecoderFunc
}

type decoderRegistry struct {
	registeredDecoders map[string]MessageDecoderFunc
}

func NewDecoderRegistry() DecoderRegistry {

	Decoders := map[string]MessageDecoderFunc{
		"tem_lab_14ns": SenlabTBasicDecoder,
	}

	return &decoderRegistry{
		registeredDecoders: Decoders,
	}
}

func (c *decoderRegistry) DesignateDecoders(ctx context.Context, sensorType string) MessageDecoderFunc {

	decoder, exist := c.registeredDecoders[sensorType]
	if exist {
		return decoder
	}

	return DefaultDecoder
}

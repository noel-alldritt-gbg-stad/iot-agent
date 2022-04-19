package decoder

import (
	"context"
)

type DecoderRegistry interface {
	GetDecodersForSensorType(ctx context.Context, sensorType string) MessageDecoderFunc
}

type decoderRegistry struct {
	registeredDecoders map[string]MessageDecoderFunc
}

func NewDecoderRegistry() DecoderRegistry {

	Decoders := map[string]MessageDecoderFunc{
		"tem_lab_14ns": SenlabTBasicDecoder,
		"Elsys_Codec": ElsysDecoder,
	}

	return &decoderRegistry{
		registeredDecoders: Decoders,
	}
}

func (c *decoderRegistry) GetDecodersForSensorType(ctx context.Context, sensorType string) MessageDecoderFunc {

	if d, ok := c.registeredDecoders[sensorType]; ok {
		return d
	}

	return DefaultDecoder
}

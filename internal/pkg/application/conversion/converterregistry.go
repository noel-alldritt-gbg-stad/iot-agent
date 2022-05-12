package conversion

import (
	"context"
)

type ConverterRegistry interface {
	DesignateConverters(ctx context.Context, types []string) []MessageConverterFunc
}

type converterRegistry struct {
	registeredConverters map[string]MessageConverterFunc
}

func NewConverterRegistry() ConverterRegistry {

	converters := map[string]MessageConverterFunc{
		"urn:oma:lwm2m:ext:3303": Temperature,
		"urn:oma:lwm2m:ext:3428": AirQuality,
		"urn:oma:lwm2m:ext:3302": Presence,
	}

	return &converterRegistry{
		registeredConverters: converters,
	}
}

func (c *converterRegistry) DesignateConverters(ctx context.Context, types []string) []MessageConverterFunc {
	converters := []MessageConverterFunc{}

	for _, t := range types {
		mc, exist := c.registeredConverters[t]
		if exist {
			converters = append(converters, mc)
		}
	}

	return converters
}

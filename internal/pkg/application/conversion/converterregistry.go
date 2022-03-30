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
	}

	return &converterRegistry{
		registeredConverters: converters,
	}
}

// bestäm vilken converter från en lista av converters som ska användas till ett visst meddelande
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

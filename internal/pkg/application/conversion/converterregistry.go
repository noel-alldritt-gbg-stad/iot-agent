package conversion

import (
	"context"
)

type ConverterRegistry interface {
	DesignateConverters(ctx context.Context, types []string) []MessageConverter
}

type converterRegistry struct {
	registeredConverters []map[string]MessageConverter
}

func NewConverterRegistry() ConverterRegistry {
	return &converterRegistry{
		registeredConverters: []map[string]MessageConverter{
			{
				"temperature": &msgConverter{
					Type: "urn:oma:lwm2m:ext:3303",
				},
			},
			{
				"presence": &msgConverter{
					Type: "presence", //this is just here because for now...
				},
			},
		},
	}
}

// bestäm vilken converter från en lista av converters som ska användas till ett visst meddelande
func (c *converterRegistry) DesignateConverters(ctx context.Context, types []string) []MessageConverter {
	converters := []MessageConverter{}

	for i, t := range types {
		mc, exist := c.registeredConverters[i][t]
		if exist {
			converters = append(converters, mc)
		}
	}

	return converters
}

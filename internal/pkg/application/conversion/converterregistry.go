package conversion

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
)

type ConverterRegistry interface {
	DesignateConverters(ctx context.Context, types []string) []MessageConverterFunc
}

type converterRegistry struct {
	registeredConverters map[string]MessageConverterFunc
}

func NewConverterRegistry() ConverterRegistry {

	var f MessageConverterFunc = func(ctx context.Context, log zerolog.Logger, internalID string, msg []byte) (*InternalMessage, error) {
		dm := struct {
			Object struct {
				Temperature float64 `json:"externalTemperature"`
			} `json:"object"`
		}{}

		err := json.Unmarshal(msg, &dm)
		if err != nil {
			return nil, err
		}

		payload := &InternalMessage{
			InternalID:  internalID,
			Type:        "urn:oma:lwm2m:ext:3303",
			SensorValue: dm.Object.Temperature,
		}
		return payload, nil
	}

	converters := map[string]MessageConverterFunc{
		"urn:oma:lwm2m:ext:3303": f,
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

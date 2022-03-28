package conversion

import (
	"context"

	"github.com/diwise/iot-agent/internal/pkg/domain"
)

type ConverterRegistry interface {
	DesignateConverters(context.Context, domain.Result) []MessageConverter
}

type converterRegistry struct {
}

func NewConverterRegistry() ConverterRegistry {
	return &converterRegistry{}
}

// bestämt vilken converter som ska användas till ett visst meddelande

func (c *converterRegistry) DesignateConverters(context.Context, domain.Result) []MessageConverter {
	return nil
}

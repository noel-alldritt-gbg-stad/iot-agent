package conversion

import "github.com/diwise/iot-agent/internal/pkg/domain"

type ConverterRegistry interface {
	DesignateConverters(domain.Result) []MessageConverter
}

type converterRegistry struct {
}

func NewConverterRegistry() ConverterRegistry {
	cr := &converterRegistry{}

	return cr
}

// bestämt vilken converter som ska användas till ett visst meddelande

func (c *converterRegistry) DesignateConverters(domain.Result) []MessageConverter {
	return nil
}

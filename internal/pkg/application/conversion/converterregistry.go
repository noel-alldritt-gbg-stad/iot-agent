package conversion

import "github.com/diwise/iot-agent/internal/pkg/domain"

type ConverterRegistry interface {
	DesignateConverters(domain.Result) []MessageConverter
}

type converterRegistry struct {
}

func NewConverterRegistry() ConverterRegistry {
	return &converterRegistry{}
}

// bestämt vilken converter från en lista av converters, som ska användas till ett visst meddelande

func (c *converterRegistry) DesignateConverters(domain.Result) []MessageConverter {
	//converters used are decided on by data format (is it from LoRa/CoAP) and type of measurements return

	return []MessageConverter{}
}

var RegisteredConverters []MessageConverter = []MessageConverter{
	&msgConverter{
		Type: "water",
	},
}

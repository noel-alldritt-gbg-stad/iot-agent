package application

type ConverterRegistry interface {
	Designate()
}

type converterRegistry struct {
}

func NewConverterRegistry() ConverterRegistry {
	cr := &converterRegistry{}

	return cr
}

//bestämt vilken converter som ska användas till ett visst meddelande

func (c *converterRegistry) Designate() {

}

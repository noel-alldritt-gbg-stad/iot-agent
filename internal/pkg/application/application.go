package application

type Application interface {
}

type iotAgentApp struct {
}

func NewApplication() Application {
	app := &iotAgentApp{}

	return app
}

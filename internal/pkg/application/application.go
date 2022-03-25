package application

type IoTAgent interface {
	NewMessage(msg []byte) error
}

type iotAgent struct {
	mp MessageProcessor
}

func NewIoTAgent(mp MessageProcessor) IoTAgent {
	app := &iotAgent{
		mp: mp,
	}

	return app
}

// this function is likely to be renamed
func (a *iotAgent) NewMessage(msg []byte) error {
	err := a.mp.ProcessMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

package domain

type DeviceManagementClient interface {
	FindDeviceFromDevEUI() (Result, error)
}

type devManagementClient struct {
}

func NewDeviceManagementClient() DeviceManagementClient {
	dmc := &devManagementClient{}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI() (Result, error) {
	// this will be a request to diff service. return internal id, type
	return Result{}, nil
}

type Result struct {
	InternalID string
	Types      []string
}

package domain

type DeviceManagementClient interface {
	FindDeviceFromDevEUI() error
}

type devManagementClient struct {
}

func NewDeviceManagementClient() DeviceManagementClient {
	dmc := &devManagementClient{}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI() error {
	// this will be a request to diff service. return internal id, type, and payload?
	return nil
}

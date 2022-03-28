package domain

import "context"

type DeviceManagementClient interface {
	FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error)
}

type devManagementClient struct {
}

func NewDeviceManagementClient() DeviceManagementClient {
	dmc := &devManagementClient{}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error) {
	// this will be a http request to diff service.

	return Result{}, nil
}

type Result struct {
	InternalID string
	Types      []string
}

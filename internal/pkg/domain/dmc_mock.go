// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package domain

import (
	"sync"
)

// Ensure, that DeviceManagementClientMock does implement DeviceManagementClient.
// If this is not the case, regenerate this file with moq.
var _ DeviceManagementClient = &DeviceManagementClientMock{}

// DeviceManagementClientMock is a mock implementation of DeviceManagementClient.
//
// 	func TestSomethingThatUsesDeviceManagementClient(t *testing.T) {
//
// 		// make and configure a mocked DeviceManagementClient
// 		mockedDeviceManagementClient := &DeviceManagementClientMock{
// 			FindDeviceFromDevEUIFunc: func(devEUI string) (Result, error) {
// 				panic("mock out the FindDeviceFromDevEUI method")
// 			},
// 		}
//
// 		// use mockedDeviceManagementClient in code that requires DeviceManagementClient
// 		// and then make assertions.
//
// 	}
type DeviceManagementClientMock struct {
	// FindDeviceFromDevEUIFunc mocks the FindDeviceFromDevEUI method.
	FindDeviceFromDevEUIFunc func(devEUI string) (Result, error)

	// calls tracks calls to the methods.
	calls struct {
		// FindDeviceFromDevEUI holds details about calls to the FindDeviceFromDevEUI method.
		FindDeviceFromDevEUI []struct {
			// DevEUI is the devEUI argument value.
			DevEUI string
		}
	}
	lockFindDeviceFromDevEUI sync.RWMutex
}

// FindDeviceFromDevEUI calls FindDeviceFromDevEUIFunc.
func (mock *DeviceManagementClientMock) FindDeviceFromDevEUI(devEUI string) (Result, error) {
	if mock.FindDeviceFromDevEUIFunc == nil {
		panic("DeviceManagementClientMock.FindDeviceFromDevEUIFunc: method is nil but DeviceManagementClient.FindDeviceFromDevEUI was just called")
	}
	callInfo := struct {
		DevEUI string
	}{
		DevEUI: devEUI,
	}
	mock.lockFindDeviceFromDevEUI.Lock()
	mock.calls.FindDeviceFromDevEUI = append(mock.calls.FindDeviceFromDevEUI, callInfo)
	mock.lockFindDeviceFromDevEUI.Unlock()
	return mock.FindDeviceFromDevEUIFunc(devEUI)
}

// FindDeviceFromDevEUICalls gets all the calls that were made to FindDeviceFromDevEUI.
// Check the length with:
//     len(mockedDeviceManagementClient.FindDeviceFromDevEUICalls())
func (mock *DeviceManagementClientMock) FindDeviceFromDevEUICalls() []struct {
	DevEUI string
} {
	var calls []struct {
		DevEUI string
	}
	mock.lockFindDeviceFromDevEUI.RLock()
	calls = mock.calls.FindDeviceFromDevEUI
	mock.lockFindDeviceFromDevEUI.RUnlock()
	return calls
}

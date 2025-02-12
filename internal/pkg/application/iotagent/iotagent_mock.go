// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package iotagent

import (
	"context"
	"sync"
)

// Ensure, that IoTAgentMock does implement IoTAgent.
// If this is not the case, regenerate this file with moq.
var _ IoTAgent = &IoTAgentMock{}

// IoTAgentMock is a mock implementation of IoTAgent.
//
//	func TestSomethingThatUsesIoTAgent(t *testing.T) {
//
//		// make and configure a mocked IoTAgent
//		mockedIoTAgent := &IoTAgentMock{
//			MessageReceivedFunc: func(ctx context.Context, msg []byte) error {
//				panic("mock out the MessageReceived method")
//			},
//		}
//
//		// use mockedIoTAgent in code that requires IoTAgent
//		// and then make assertions.
//
//	}
type IoTAgentMock struct {
	// MessageReceivedFunc mocks the MessageReceived method.
	MessageReceivedFunc func(ctx context.Context, msg []byte) error

	// calls tracks calls to the methods.
	calls struct {
		// MessageReceived holds details about calls to the MessageReceived method.
		MessageReceived []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Msg is the msg argument value.
			Msg []byte
		}
	}
	lockMessageReceived sync.RWMutex
}

// MessageReceived calls MessageReceivedFunc.
func (mock *IoTAgentMock) MessageReceived(ctx context.Context, msg []byte) error {
	if mock.MessageReceivedFunc == nil {
		panic("IoTAgentMock.MessageReceivedFunc: method is nil but IoTAgent.MessageReceived was just called")
	}
	callInfo := struct {
		Ctx context.Context
		Msg []byte
	}{
		Ctx: ctx,
		Msg: msg,
	}
	mock.lockMessageReceived.Lock()
	mock.calls.MessageReceived = append(mock.calls.MessageReceived, callInfo)
	mock.lockMessageReceived.Unlock()
	return mock.MessageReceivedFunc(ctx, msg)
}

// MessageReceivedCalls gets all the calls that were made to MessageReceived.
// Check the length with:
//
//	len(mockedIoTAgent.MessageReceivedCalls())
func (mock *IoTAgentMock) MessageReceivedCalls() []struct {
	Ctx context.Context
	Msg []byte
} {
	var calls []struct {
		Ctx context.Context
		Msg []byte
	}
	mock.lockMessageReceived.RLock()
	calls = mock.calls.MessageReceived
	mock.lockMessageReceived.RUnlock()
	return calls
}

package events

import "time"

type StatusMessage struct {
	DeviceID  string `json:"deviceID"`
	Timestamp string `json:"timestamp"`
}

func NewStatusMessage(deviceID string) *StatusMessage {
	msg := &StatusMessage{
		DeviceID:  deviceID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	return msg
}

func (m *StatusMessage) ContentType() string {
	return "application/json"
}

func (m *StatusMessage) TopicName() string {
	return "Status"
}

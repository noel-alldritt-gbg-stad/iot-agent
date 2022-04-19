package decoder

import "context"

type Payload struct {
	DevEUI       string        `json:"devEUI"`
	DeviceName   string        `json:"deviceName,omitempty"`
	FPort        int           `json:"fPort,omitempty"`
	Latitude     float64       `json:"latitude,omitempty"`
	Longitude    float64       `json:"longitude,omitempty"`
	Rssi         int           `json:"rssi,omitempty"`
	SensorType   string        `json:"sensorType,omitempty"`
	Timestamp    string        `json:"timestamp,omitempty"`
	Type         string        `json:"type,omitempty"`
	Error        string        `json:"error,omitempty"`
	Measurements []interface{} `json:"measurements"`
}

type MessageDecoderFunc func(context.Context, []byte, func(context.Context, []byte) error) error

func DefaultDecoder(ctx context.Context, msg []byte, fn func(context.Context, []byte) error) error {
	err := fn(ctx, msg)
	return err
}

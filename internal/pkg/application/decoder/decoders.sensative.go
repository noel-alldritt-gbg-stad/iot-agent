package decoder

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
)

func SensativeDecoder(ctx context.Context, msg []byte, fn func(context.Context, []byte) error) error {

	dm := []struct {
		DevEUI     string  `json:"devEUI"`
		FPort      string  `json:"fPort,omitempty"`
		Latitude   float64 `json:"latitude,omitempty"`
		Longitude  float64 `json:"longitude,omitempty"`
		Rssi       string  `json:"rssi,omitempty"`
		SensorType string  `json:"sensorType,omitempty"`
		Timestamp  string  `json:"timestamp,omitempty"`
		Payload    string  `json:"payload"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return err
	}

	for _, d := range dm {

		b, err := hex.DecodeString(d.Payload)
		if err != nil {
			return err
		}

		if len(b) < 4 {
			return errors.New("payload too short")
		}

		pp := &Payload{
			DevEUI:       d.DevEUI,
			FPort:        d.FPort,
			Latitude:     d.Latitude,
			Longitude:    d.Longitude,
			Rssi:         d.Rssi,
			SensorType:   d.SensorType,
			Timestamp:    d.Timestamp,
			Measurements: []any{},
		}

		err = decodeSensativeMeasurements(b, func(m Measurement) {
			pp.Measurements = append(pp.Measurements, m)
		})
		if err != nil {
			return err
		}

		r, err := json.Marshal(&pp)
		if err != nil {
			return nil
		}

		err = fn(ctx, r)
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeSensativeMeasurements(payload []byte, callback func(m Measurement)) error {
	temp := struct {
		Value float32 `json:"temperature"`
	}{42.0}

	callback(temp)

	return nil
}

type Measurement interface{}

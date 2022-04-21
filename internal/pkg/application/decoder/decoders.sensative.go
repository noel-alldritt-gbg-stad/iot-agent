package decoder

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
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

		if len(b) < 2 {
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

	pos := 2

	for pos < len(payload) {
		channel := payload[pos] & 0x7F
		pos = pos + 1
		size := 1

		switch channel {
		case 1: // battery
			callback(struct {
				Value int `json:"battery_level"`
			}{int(payload[pos])})
		case 2: // temp report
			size = 2
			// TODO: Handle sub zero readings
			callback(struct {
				Value float64 `json:"temperature"`
			}{float64(binary.BigEndian.Uint16(payload[pos:pos+2])) / 10})
		case 4: // average temp report
			size = 2
		case 6: // humidity report
			callback(struct {
				Value int `json:"humidity"`
			}{int(payload[pos] / 2)})
		case 7: // lux report
			size = 2
		case 8: // lux2 report
			size = 2
		case 9: // door report
			callback(struct {
				Value bool `json:"door_report"`
			}{payload[pos] != 0})
		case 10: // door alarm
			callback(struct {
				Value bool `json:"door_alarm"`
			}{payload[pos] != 0})
		default:
			fmt.Printf("unknown channel %d\n", channel)
			size = 20
		}

		pos = pos + size
	}

	return nil
}

type Measurement interface{}

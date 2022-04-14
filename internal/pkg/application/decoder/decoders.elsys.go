package decoder

import (
	"context"
	"encoding/json"
)

func ElsysDecoder(ctx context.Context, msg []byte, fn func(context.Context, []byte) error) error {

	d := struct {
		DevEUI     string `json:"devEUI"`
		FPort      int    `json:"fPort,omitempty"`
		SensorType string `json:"deviceProfileName,omitempty"`
		Data       string `json:"data"`
		Object     struct {
			Temperature         *float32 `json:"temperature,omitempty"`
			ExternalTemperature *float32 `json:"externalTemperature,omitempty"`
			Vdd                 *int     `json:"vdd,omitempty"`
		} `json:"object"`
	}{}

	err := json.Unmarshal(msg, &d)
	if err != nil {
		return err
	}

	pp := &Payload{
		DevEUI:     d.DevEUI,
		FPort:      d.FPort,
		SensorType: d.SensorType,
	}

	if d.Object.Temperature != nil {
		temp := struct {
			Temperature float32
		}{
			*d.Object.Temperature,
		}
		pp.Measurements = append(pp.Measurements, temp)
	}

	if d.Object.ExternalTemperature != nil {
		temp := struct {
			Temperature float32 `json:"temperature"`
		}{
			*d.Object.ExternalTemperature,
		}
		pp.Measurements = append(pp.Measurements, temp)
	}

	if d.Object.Vdd != nil {
		bat := struct {
			BatteryCurrentLevel int `json:"battery_current_level"`
		}{
			*d.Object.Vdd,
		}
		pp.Measurements = append(pp.Measurements, bat)
	}

	r, err := json.Marshal(&pp)
	if err != nil {
		return nil
	}

	err = fn(ctx, r)
	if err != nil {
		return err
	}	

	return nil
}



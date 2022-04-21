package decoder

import (
	"context"
	"encoding/json"
	"fmt"
)

func ElsysDecoder(ctx context.Context, msg []byte, fn func(context.Context, []byte) error) error {

	d := struct {
		DevEUI     string `json:"devEUI"`
		FPort      string `json:"fPort"`
		SensorType string `json:"deviceProfileName"`
		Data       string `json:"data"`
		Object     struct {
			Temperature         *float32 `json:"temperature,omitempty"`
			ExternalTemperature *float32 `json:"externalTemperature,omitempty"`
			Vdd                 *int     `json:"vdd,omitempty"`
		} `json:"object"`
	}{}

	err := json.Unmarshal(msg, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal elsys payload: %s", err.Error())
	}

	pp := &Payload{
		DevEUI:     d.DevEUI,
		FPort:      d.FPort,
		SensorType: d.SensorType,
	}

	if d.Object.Temperature != nil {
		temp := struct {
			Temperature float32 `json:"temperature"`
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
			BatteryLevel int `json:"battery_level"`
		}{
			*d.Object.Vdd, // TODO: Adjust for max VDD
		}
		pp.Measurements = append(pp.Measurements, bat)
	}

	r, err := json.Marshal(&pp)
	if err != nil {
		return err
	}

	err = fn(ctx, r)
	if err != nil {
		return err
	}

	return nil
}

package decoder

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func WatermeteringDecoder(ctx context.Context, msg []byte, fn func(context.Context, Payload) error) error {

	d := struct {
		DevEUI     string `json:"devEUI"`
		FPort      int    `json:"fPort"`
		SensorType string `json:"deviceProfileName"`
		Data       string `json:"data"`
		Object     struct {
			CurrentDateTime string   `json:"curDateTime"`
			CurrentVolume   *float32 `json:"curVol,omitempty"`
		} `json:"object"`
	}{}

	err := json.Unmarshal(msg, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal elsys payload: %s", err.Error())
	}

	pp := &Payload{
		DevEUI:     d.DevEUI,
		FPort:      strconv.Itoa(d.FPort),
		SensorType: d.SensorType,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	if d.Object.CurrentDateTime != "" {
		curDateTime := struct {
			CurrentDateTime string `json:"curDateTime"`
		}{
			d.Object.CurrentDateTime,
		}
		pp.Measurements = append(pp.Measurements, curDateTime)
	}

	if d.Object.CurrentVolume != nil {
		curVol := struct {
			CurrentVolume float32 `json:"curVol"`
		}{
			*d.Object.CurrentVolume,
		}
		pp.Measurements = append(pp.Measurements, curVol)
	}

	err = fn(ctx, *pp)
	if err != nil {
		return err
	}

	return nil
}

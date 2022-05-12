package decoder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func PresenceDecoder(ctx context.Context, msg []byte, fn func(context.Context, Payload) error) error {
	d := struct {
		DevEUI string `json:"devEUI"`
		Data   string `json:"data"`
		Object struct {
			Presence *bool `json:"present,omitempty"`
		} `json:"object,omitempty"`
		ObjectJSON struct {
			BuildID struct {
				Id       *int64 `json:"id,omitempty"`
				Modified *bool  `json:"modified,omitempty"`
			} `json:"buildId,omitempty"`
			CloseProximityAlarm struct {
				Value *bool `json:"value,omitempty"`
			} `json:"closeProximityAlarm"`
		} `json:"objectJSON,omitempty"`
	}{}

	err := json.Unmarshal(msg, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal presence payload: %s", err.Error())
	}

	payload := &Payload{
		DevEUI:    d.DevEUI,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if d.Object.Presence != nil {
		present := struct {
			Presence bool `json:"present"`
		}{
			*d.Object.Presence,
		}
		payload.Measurements = append(payload.Measurements, present)
	}
	if d.ObjectJSON.CloseProximityAlarm.Value != nil {
		present := struct {
			Presence bool `json:"present"`
		}{
			*d.ObjectJSON.CloseProximityAlarm.Value,
		}
		payload.Measurements = append(payload.Measurements, present)
	}

	return fn(ctx, *payload)
}

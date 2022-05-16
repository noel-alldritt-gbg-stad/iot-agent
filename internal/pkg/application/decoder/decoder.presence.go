package decoder

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func PresenceDecoder(ctx context.Context, msg []byte, fn func(context.Context, Payload) error) error {
	d := struct {
		DevEUI            string `json:"devEUI"`
		Data              string `json:"data"`
		DeviceProfileName string `json:"deviceProfileName"`
		Object            struct {
			Presence struct {
				Value *bool `json:"value"`
			} `json:"closeProximityAlarm,omitempty"`
		} `json:"object,omitempty"`
	}{}

	err := json.Unmarshal(msg, &d)
	if err != nil {
		return fmt.Errorf("failed to unmarshal presence payload: %s", err.Error())
	}

	payload := &Payload{
		DevEUI:    d.DevEUI,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if d.Object.Presence.Value != nil {
		present := struct {
			Presence bool `json:"present"`
		}{
			*d.Object.Presence.Value,
		}
		payload.Measurements = append(payload.Measurements, present)
	}

	return fn(ctx, *payload)
}

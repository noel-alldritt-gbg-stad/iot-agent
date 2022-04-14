package conversion

import (
	"context"
	"encoding/json"
	"fmt"
)

type MessageConverterFunc func(ctx context.Context, internalID string, msg []byte) (*InternalMessage, error)

func Temperature(ctx context.Context, internalID string, msg []byte) (*InternalMessage, error) {
	dm := struct {
		Measurements []struct {			
			Temp    *float64 `json:"temperature"`
		} `json:"measurements"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return nil, err
	}

	payload := &InternalMessage{
		InternalID: internalID,
		Type:       "urn:oma:lwm2m:ext:3303",
	}

	//TODO: range and call func?
 	if dm.Measurements[0].Temp != nil {
		payload.SensorValue = *dm.Measurements[0].Temp
	} else {
		return nil, fmt.Errorf("no temperature value found in payload")
	}

	return payload, nil
}

type InternalMessage struct {
	InternalID  string  `json:"internalID"`
	Type        string  `json:"type"`
	SensorValue float64 `json:"sensorValue"`
}

func (im InternalMessage) ContentType() string {
	return "application/json" // TODO: Decide a proper content type here
}

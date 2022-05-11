package conversion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/diwise/iot-agent/internal/pkg/application/decoder"
	"github.com/farshidtz/senml/v2"
)

type MessageConverterFunc func(ctx context.Context, internalID string, payload decoder.Payload) (senml.Pack, error)

func Temperature(ctx context.Context, deviceID string, payload decoder.Payload) (senml.Pack, error) {
	dm := struct {
		Timestamp    string `json:"timestamp"`
		Measurements []struct {
			Temp *float64 `json:"temperature"`
		} `json:"measurements"`
	}{}

	if err := convertPayloadToStruct(payload, &dm); err != nil {
		return nil, fmt.Errorf("failed to convert payload: %s", err.Error())
	}

	baseTime, err := parseTime(dm.Timestamp)
	if err != nil {
		return nil, err
	}

	var pack senml.Pack
	pack = append(pack, senml.Record{
		BaseName:    "urn:oma:lwm2m:ext:3303",
		BaseTime:    baseTime,
		Name:        "0",
		StringValue: deviceID,
	})

	for _, m := range dm.Measurements {
		if m.Temp != nil {
			rec := senml.Record{
				Name:  "Temperature",
				Value: m.Temp,
			}

			pack = append(pack, rec)
		}
	}

	return pack, nil
}

func AirQuality(ctx context.Context, deviceID string, payload decoder.Payload) (senml.Pack, error) {
	dm := struct {
		Timestamp    string `json:"timestamp"`
		Measurements []struct {
			CO2 *int `json:"co2"`
		} `json:"measurements"`
	}{}

	if err := convertPayloadToStruct(payload, &dm); err != nil {
		return nil, fmt.Errorf("failed to convert payload: %s", err.Error())
	}

	baseTime, err := parseTime(dm.Timestamp)
	if err != nil {
		return nil, err
	}

	var pack senml.Pack
	pack = append(pack, senml.Record{
		BaseName:    "urn:oma:lwm2m:ext:3428",
		BaseTime:    baseTime,
		Name:        "0",
		StringValue: deviceID,
	})

	for _, m := range dm.Measurements {
		if m.CO2 != nil {
			co2 := float64(*m.CO2)
			rec := senml.Record{
				Name:  "CO2",
				Value: &co2,
			}

			pack = append(pack, rec)
		}
	}

	return pack, nil
}

func parseTime(t string) (float64, error) {
	baseTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return 0, fmt.Errorf("unable to parse time %s as RFC3339, %s", t, err.Error())
	}

	return float64(baseTime.Unix()), nil
}

func convertPayloadToStruct(p decoder.Payload, v any) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	return nil
}

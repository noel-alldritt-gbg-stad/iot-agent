package conversion

import (
	"encoding/json"
)

type MessageConverter interface {
	ConvertPayload(msg []byte) (InternalMessageFormat, error)
}

//konvertera payload till internt format

type msgConverter struct {
	Type string //determines what type of data we're converting, i.e. water or air temperature etc.
}

func (mc *msgConverter) ConvertPayload(msg []byte) (InternalMessageFormat, error) {
	imf := InternalMessageFormat{}

	err := json.Unmarshal(msg, &imf)
	if err != nil {
		return imf, err
	}

	return imf, nil
}

type InternalMessageFormat struct {
	InternalID string
	Types      []string
	Longitude  float64
	Latitude   float64
	Value      float64
}

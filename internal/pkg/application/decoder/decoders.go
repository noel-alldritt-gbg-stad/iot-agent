package decoder

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type MessageDecoderFunc func(ctx context.Context, msg []byte) ([]byte, error)

func DefaultDecoder(ctx context.Context, msg []byte) ([]byte, error) {
	return msg, nil
}

func SenlabTBasicDecoder(ctx context.Context, msg []byte) ([]byte, error) {

	dm := []struct {
		DevEUI    string `json:"devEui"`
		Payload   string `json:"payload"`
		Timestamp string `json:"timestamp"`
	}{}

	err := json.Unmarshal(msg, &dm)
	if err != nil {
		return nil, err
	}

	//TODO: range?
	b, err := hex.DecodeString(dm[0].Payload)
	if err != nil {
		return nil, err // TODO: do something better?!
	}

	var temp int16
	err = binary.Read(bytes.NewReader(b[len(b)-2:]), binary.BigEndian, &temp)
	if err != nil {
		return nil, err // TODO: do something better?!
	}

	id := int(b[0])
	battery := (int(b[1]) / 254) * 100
	temperature := float32(temp) / 16.0

	// these values must be ignored since they are sensor reading errors
	if temperature == -46.75 || temperature == 85 {
		return nil, errors.New("sensor reading error")
	}

	result := struct {
		DevEUI       string
		Id           int
		BatteryLevel int
		Temperature  float32
		Timestamp    string
	}{
		dm[0].DevEUI,
		id,
		battery,
		temperature,
		dm[0].Timestamp,
	}

	reqBodyBytes := new(bytes.Buffer)
	err = json.NewEncoder(reqBodyBytes).Encode(result)
	if err != nil {
		return nil, err // TODO: do something better?!
	}

	return reqBodyBytes.Bytes(), nil
}

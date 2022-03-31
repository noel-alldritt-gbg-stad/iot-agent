package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

type DeviceManagementClient interface {
	FindDeviceFromDevEUI(ctx context.Context, devEUI string) (*Result, error)
}

type devManagementClient struct {
	url string
	log zerolog.Logger
}

var tracer = otel.Tracer("dmc-client")

func NewDeviceManagementClient(devMgmtUrl string, log zerolog.Logger) DeviceManagementClient {
	dmc := &devManagementClient{
		url: devMgmtUrl,
		log: log,
	}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI(ctx context.Context, devEUI string) (*Result, error) {
	var err error
	ctx, span := tracer.Start(ctx, "find-device")
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	dmc.log.Info().Msgf("looking up internal id and types for devEUI %s", devEUI)

	resp, err := http.Get(dmc.url + "/api/v0/devices/" + devEUI)
	if err != nil {
		dmc.log.Error().Msgf("failed to retrieve device information from devEUI: %s", err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		dmc.log.Error().Msgf("request failed with status code %d", resp.StatusCode)
		return nil, fmt.Errorf("request failed, no device found")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dmc.log.Error().Msgf("failed to read response body: %s", err.Error())
		return nil, err
	}

	result := &Result{}

	err = json.Unmarshal(respBody, result)
	if err != nil {
		dmc.log.Error().Msgf("failed to unmarshal response body: %s", err.Error())
		return nil, err
	}

	return result, nil
}

type Result struct {
	InternalID string   `json:"id"`
	Types      []string `json:"types"`
}

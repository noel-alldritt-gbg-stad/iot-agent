package domain

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
)

type DeviceManagementClient interface {
	FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error)
}

type devManagementClient struct {
	url string
	log zerolog.Logger
}

func NewDeviceManagementClient(dmcurl string, log zerolog.Logger) DeviceManagementClient {
	dmc := &devManagementClient{
		url: dmcurl,
		log: log,
	}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error) {
	// this will be a http request to diff service.
	result := Result{}

	resp, err := http.Get(dmc.url + "/" + devEUI)
	if resp.StatusCode != http.StatusOK {
		dmc.log.Error().Msgf("request failed with status code %d", resp.StatusCode)
		return result, nil
	}
	if err != nil {
		dmc.log.Error().Msgf("failed to retrieve device information from devEUI: %s", err.Error())
		return result, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		dmc.log.Error().Msgf("failed to read response body: %s", err.Error())
		return result, err
	}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		dmc.log.Error().Msgf("failed to unmarshal response body: %s", err.Error())
		return result, err
	}

	return result, nil
}

type Result struct {
	InternalID string   `json:"internalID"`
	Types      []string `json:"types"`
}

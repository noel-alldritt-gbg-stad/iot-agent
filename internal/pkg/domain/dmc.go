package domain

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

type DeviceManagementClient interface {
	FindDeviceFromDevEUI(ctx context.Context, devEUI string) (*Result, error)
}

type devManagementClient struct {
	url string
}

var tracer = otel.Tracer("dmc-client")

func NewDeviceManagementClient(devMgmtUrl string) DeviceManagementClient {
	dmc := &devManagementClient{
		url: devMgmtUrl,
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

	log := logging.GetFromContext(ctx)

	log.Info().Msgf("looking up internal id and types for devEUI %s", devEUI)

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	url := dmc.url + "/api/v0/devices?devEUI=" + devEUI

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to create http request")
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Error().Msgf("failed to retrieve device information from devEUI: %s", err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("request failed with status code %d", resp.StatusCode)
		return nil, fmt.Errorf("request failed, no device found")
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Msgf("failed to read response body: %s", err.Error())
		return nil, err
	}

	result := []Result{}

	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Error().Msgf("failed to unmarshal response body: %s", err.Error())
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("device management returned an empty list of devices")
	}

	device := result[0]
	return &device, nil
}

type Result struct {
	InternalID string   `json:"id"`
	SensorType string   `json:"sensorType"`
	Types      []string `json:"types"`
}

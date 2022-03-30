package domain

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
)

type DeviceManagementClient interface {
	FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error)
}

type devManagementClient struct {
	url string
	log zerolog.Logger
}

var tracer = otel.Tracer("dmc-client")

func NewDeviceManagementClient(dmcurl string, log zerolog.Logger) DeviceManagementClient {
	dmc := &devManagementClient{
		url: dmcurl,
		log: log,
	}
	return dmc
}

func (dmc *devManagementClient) FindDeviceFromDevEUI(ctx context.Context, devEUI string) (Result, error) {

	ctx, span := tracer.Start(ctx, "find-device")
	defer span.End()

	dmc.log.Info().Msgf("looking up internal id and types for devEUI %s", devEUI)

	// this will be a http request to diff service.
	result := Result{
		InternalID: fmt.Sprintf(
			"internalID:%s", devEUI),
		Types: []string{"urn:oma:lwm2m:ext:3303"},
	}

	/*resp, err := http.Get(dmc.url + "/" + devEUI)
	if err != nil {
		dmc.log.Error().Msgf("failed to retrieve device information from devEUI: %s", err.Error())
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		dmc.log.Error().Msgf("request failed with status code %d", resp.StatusCode)
		return result, nil
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
	}*/

	return result, nil
}

type Result struct {
	InternalID string   `json:"internalID"`
	Types      []string `json:"types"`
}

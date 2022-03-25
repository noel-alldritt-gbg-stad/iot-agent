package api

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diwise/iot-agent/internal/pkg/application"
	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

func TestHealthEndpointReturns204StatusNoContent(t *testing.T) {
	is, a, _ := testSetup(t)

	server := httptest.NewServer(a.r)
	defer server.Close()

	resp, _ := testRequest(is, server, http.MethodGet, "/health", nil)
	is.Equal(resp.StatusCode, http.StatusNoContent)
}

func TestXxx(t *testing.T) {
	is, api, app := testSetup(t)

	server := httptest.NewServer(api.r)
	defer server.Close()

	resp, _ := testRequest(is, server, http.MethodPost, "/newmsg", bytes.NewBuffer([]byte(msgfromMQTT)))
	is.Equal(resp.StatusCode, http.StatusCreated)
	is.Equal(len(app.MessageReceivedCalls()), 1)
}

func testSetup(t *testing.T) (*is.I, *api, *application.IoTAgentMock) {
	is := is.New(t)
	r := chi.NewRouter()

	app := &application.IoTAgentMock{
		MessageReceivedFunc: func(msg []byte) error {
			return nil
		},
	}

	a := newAPI(r, app)

	return is, a, app
}

func testRequest(is *is.I, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, _ := http.NewRequest(method, ts.URL+path, body)
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp, string(respBody)
}

const msgfromMQTT string = `{{"applicationID":"53","applicationName":"POC-SC-IT","deviceName":"Elsys_EMS_3","deviceProfileName":"Elsys_codec","deviceProfileID":"xxxxxxx","devEUI":"cnajncksla","rxInfo":[{"gatewayID":"fcc23dfffe0a752b","uplinkID":"46ae2624-0978-45a1-bf11-f8fa9f91763c","name":"SN-LGW-001","time":"2022-03-25T09:41:57.446571673Z","rssi":-106,"loRaSNR":-14.8,"location":{"latitude":62.39466886148298,"longitude":17.34076023101807,"altitude":0}}],"txInfo":{"frequency":868300000,"dr":2},"adr":true,"fCnt":38701,"fPort":5,"data":"xxxxxxx","object":{"accMotion":0,"digital":0,"humidity":29,"pulseAbs":90,"temperature":21.9,"vdd":3596,"waterleak":0,"x":-64,"y":0,"z":-1},"tags":{"place":"pof_door"}}}`

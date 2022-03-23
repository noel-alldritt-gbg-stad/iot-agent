package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/matryer/is"
)

func TestHealthEndpointReturns204StatusNoContent(t *testing.T) {
	is, a := testSetup(t)

	server := httptest.NewServer(a.r)
	defer server.Close()

	resp, _ := testRequest(is, server, http.MethodGet, "/health", nil)
	is.Equal(resp.StatusCode, http.StatusNoContent)
}

func testSetup(t *testing.T) (*is.I, *iotAgentApi) {
	is := is.New(t)
	r := chi.NewRouter()
	a := newIotAgentApi(r)

	return is, a
}

func testRequest(is *is.I, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, _ := http.NewRequest(method, ts.URL+path, body)
	resp, _ := http.DefaultClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp, string(respBody)
}

package api

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/thermokarst/bactdb/datastore"
	"github.com/thermokarst/bactdb/models"
)

func init() {
	serveMux.Handle("/", http.StripPrefix("/api", Handler()))
}

var (
	serveMux   = http.NewServeMux()
	httpClient = http.Client{Transport: (*muxTransport)(serveMux)}
	apiClient  = models.NewClient(&httpClient)
)

func setup() {
	store = datastore.NewMockDatastore()
}

type muxTransport http.ServeMux

// RoundTrip is for testing API requests. It intercepts all requests during testing
// to serve up a local/internal response.
func (t *muxTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rw := httptest.NewRecorder()
	rw.Body = new(bytes.Buffer)
	(*http.ServeMux)(t).ServeHTTP(rw, req)
	return &http.Response{
		StatusCode:    rw.Code,
		Status:        http.StatusText(rw.Code),
		Header:        rw.HeaderMap,
		Body:          ioutil.NopCloser(rw.Body),
		ContentLength: int64(rw.Body.Len()),
		Request:       req,
	}, nil
}

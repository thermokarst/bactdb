package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/thermokarst/bactdb/datastore"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/router"
)

func init() {
	serveMux.Handle("/", http.StripPrefix("/api", Handler()))
}

var (
	serveMux   = http.NewServeMux()
	httpClient = http.Client{
		Transport: (*muxTransport)(serveMux),
	}
	apiClient = models.NewClient(&httpClient)
	testToken models.UserSession
)

func setup() {
	store = datastore.NewMockDatastore()
	SetupCerts("../keys/")
	u, _ := apiClient.URL(router.GetToken, nil, nil)
	resp, _ := httpClient.PostForm(u.String(),
		url.Values{"username": {"test_user"}, "password": {"password"}})
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&testToken); err != nil {
		panic(err)
	}
}

type muxTransport http.ServeMux

// RoundTrip is for testing API requests. It intercepts all requests during testing
// to serve up a local/internal response.
func (t *muxTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rw := httptest.NewRecorder()
	rw.Body = new(bytes.Buffer)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", testToken.Token))
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

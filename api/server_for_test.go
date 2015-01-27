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
	signKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAyXi+Q3rVv+NkBi3UWBXyylUN5SyHxQUtZvhbA6TwcF2fk3Io
Di3kipVWbfYAYijIxczACieCnadSynJpH4zNVJsRxp8DTG67Nx/K5n3TJyg5hLpa
PE+46lOS7E0O9JPT119zKCGHtHldpgjsPCXGHRXKZdFafSv8ktFhwvK5ZQO1NjH2
+NOsfhp2ubuQXL7O/45fC5wTCj0lLatdXtlTcIxJb7FUj3AsAC7TlKhtkKe+Syox
n1xNMQiYK1R24+goW44JO5uZDYP85ufgRn+DdOQF9DskmiQN9/REhH/5VQjEIYZ3
kjmKlNnjz3Jd1eOdtGALxTq3+neVawWMPBXOcQIDAQABAoIBAHMAYhKgrhxPTwwb
4ua4+JK4BCt5xLIYp3bscv9cigaJ2onOksCtP5Q/dEtmLYfaYehOXJwvO2aEWUTI
E+t3cslFjtsCb16UonbvxeDVl871LgfuW42rsBDJzcbmoY/IRhbdHB2fLhg9YtBg
rYATy8dUZejCnNVwY0bnD9e4t0zJ0lXUVy+dMvl69uNyP8f12LwvLGgCmAOWXh5p
NpGmT8/jRF9BrQvr9bhwxpV2JGsGEEyGvu+ayVR01AiyQ04kh9gZOJOVtsGa1fjx
AvgxzhkfLyAbCAgFTTUuhEbZoxXyCNBdOM0V3PXSbIbW+7gwLwXi71Czo08V050z
5SK9p2UCgYEA8JW+xIaAzYT+ZwPaJ/Ir1+WcisEI+AyCD1c5gblHrCmSwYHdKSX4
ZcX0QAcj+dzuF6SyQStoy1pIooUzadDZxXeyBoeOjGdyobqJpmaEHb9g594w2inS
AsEb4waxvrKlTuhFXnI2JbJrbMyjRBTKWZw4K/FT73IE8hQL9ivXYN8CgYEA1mFu
uLD95N/KOFtQ0JcLUelPZK7gYSdq21YHSbABeeesvabnn3hFEzDl5cXqKFJ/4Ltf
2QhpO4JGgHYNlrsfCvvivV5dRxFKfleoojL/0qlJxOqQVfulscT0XB3wUpoyP+Qr
8AdyvZwUfLWpSaYxDUB7w77U1VayP5JLuULKKq8CgYBOge8QnnullUKXRzCHXIVm
HG1q8fcFSr+eVe5UIKv8yEw1jTUoWlWmkGRWCH566NdhK8NndMzrnviY4DKY0yhd
QeP8MXwY4SENGZwVitqOAoeS4nS6nG8FqxJ4kRSrkAxVpYINgeOdhY18oYKdktM9
Trcdz9B+EI0Amf4VRNUxrQKBgQCTBXTagr9MhFF5vt44fy3LOhcxtGC7ID4/N8t9
tI/+m2yzD9DPY7rzg1hW8Rk6GAINDFOaUxNgNWLGXK/LDH8omEASoLGVuHz/Enza
5+DcBy9JNZhQ72jd9nWi6wFSlN8bRA8B6Qm+kVjXgfocQTZooS1/u9LYkEFkKZ92
6SAejwKBgH6V+dLUefC9w6N99VWpOAom9gE96imo1itTKxIB1WPdFsSDDe/CGXDG
W+l7ZiSjXaAfNF5UtuhO6yY0ob++Aa/d+l+Zo7CWn4SrmwTAp1CdRHxX3KxkqHNi
BsuYClbQh5Z9lOKn8FCNW3NyahJdYeWGhb/ZdeS0n+F6Ov4V+grc
-----END RSA PRIVATE KEY-----`)

	verifyKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyXi+Q3rVv+NkBi3UWBXy
ylUN5SyHxQUtZvhbA6TwcF2fk3IoDi3kipVWbfYAYijIxczACieCnadSynJpH4zN
VJsRxp8DTG67Nx/K5n3TJyg5hLpaPE+46lOS7E0O9JPT119zKCGHtHldpgjsPCXG
HRXKZdFafSv8ktFhwvK5ZQO1NjH2+NOsfhp2ubuQXL7O/45fC5wTCj0lLatdXtlT
cIxJb7FUj3AsAC7TlKhtkKe+Syoxn1xNMQiYK1R24+goW44JO5uZDYP85ufgRn+D
dOQF9DskmiQN9/REhH/5VQjEIYZ3kjmKlNnjz3Jd1eOdtGALxTq3+neVawWMPBXO
cQIDAQAB
-----END PUBLIC KEY-----`)

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

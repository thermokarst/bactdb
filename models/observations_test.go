package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newObservation() *Observation {
	observation := NewObservation()
	observation.Id = 1
	return observation
}

func TestObservationService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newObservation()

	var called bool
	mux.HandleFunc(urlPath(t, router.Observation, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	observation, err := client.Observations.Get(want.Id)
	if err != nil {
		t.Errorf("Observations.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(observation, want) {
		t.Errorf("Observations.Get return %+v, want %+v", observation, want)
	}
}

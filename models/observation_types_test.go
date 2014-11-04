package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newObservationType() *ObservationType {
	observation_type := NewObservationType()
	observation_type.Id = 1
	return observation_type
}

func TestObservation_TypeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newObservationType()

	var called bool
	mux.HandleFunc(urlPath(t, router.ObservationType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	observation_type, err := client.ObservationTypes.Get(want.Id)
	if err != nil {
		t.Errorf("ObservationTypes.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(observation_type, want) {
		t.Errorf("ObservationTypes.Get return %+v, want %+v", observation_type, want)
	}
}

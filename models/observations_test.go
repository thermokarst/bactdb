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

func TestObservationService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newObservation()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateObservation, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"observation_name":"Test Observation","observation_type_id":0,"created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	observation := newObservation()
	created, err := client.Observations.Create(observation)
	if err != nil {
		t.Errorf("Observations.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(observation, want) {
		t.Errorf("Observations.Create returned %+v, want %+v", observation, want)
	}
}

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
		testBody(t, r, `{"id":1,"observationName":"Test Observation","observationTypeId":0,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")

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

func TestObservationService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Observation{newObservation()}

	var called bool
	mux.HandleFunc(urlPath(t, router.Observations, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	observations, err := client.Observations.List(nil)
	if err != nil {
		t.Errorf("Observations.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(observations, want) {
		t.Errorf("Observations.List return %+v, want %+v", observations, want)
	}
}

func TestObservationService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newObservation()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateObservation, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"observationName":"Test Obs Updated","observationTypeId":0,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	observation := newObservation()
	observation.ObservationName = "Test Obs Updated"
	updated, err := client.Observations.Update(observation.Id, observation)
	if err != nil {
		t.Errorf("Observations.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestObservationService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newObservation()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteObservation, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.Observations.Delete(want.Id)
	if err != nil {
		t.Errorf("Observations.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

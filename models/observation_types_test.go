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

func TestObservationTypeService_Get(t *testing.T) {
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

func TestObservationTypeService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newObservationType()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateObservationType, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"observation_type_name":"Test Obs Type","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	observation_type := newObservationType()
	created, err := client.ObservationTypes.Create(observation_type)
	if err != nil {
		t.Errorf("ObservationTypes.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(observation_type, want) {
		t.Errorf("ObservationTypes.Create returned %+v, want %+v", observation_type, want)
	}
}

func TestObservationTypeService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*ObservationType{newObservationType()}

	var called bool
	mux.HandleFunc(urlPath(t, router.ObservationTypes, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	observation_types, err := client.ObservationTypes.List(nil)
	if err != nil {
		t.Errorf("ObservationTypes.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(observation_types, want) {
		t.Errorf("ObservationTypes.List return %+v, want %+v", observation_types, want)
	}
}

func TestObservationTypeService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newObservationType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateObservationType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"observation_type_name":"Test Obs Type Updated","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	observation_type := newObservationType()
	observation_type.ObservationTypeName = "Test Obs Type Updated"
	updated, err := client.ObservationTypes.Update(observation_type.Id, observation_type)
	if err != nil {
		t.Errorf("ObservationTypes.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

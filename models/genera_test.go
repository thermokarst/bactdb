package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func TestGeneraService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Genus{Id: 1, GenusName: "Test Genus"}

	var called bool
	mux.HandleFunc(urlPath(t, router.Genus, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	genus, err := client.Genera.Get(1)
	if err != nil {
		t.Errorf("Genera.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(genus, want) {
		t.Errorf("Genera.Get returned %+v, want %+v", genus, want)
	}
}

func TestGeneraService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := &Genus{Id: 1, GenusName: "Test Genus"}

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateGenus, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"genus_name":"Test Genus","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	genus := &Genus{Id: 1, GenusName: "Test Genus"}
	created, err := client.Genera.Create(genus)
	if err != nil {
		t.Errorf("Genera.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(genus, want) {
		t.Errorf("Genera.Create returned %+v, want %+v", genus, want)
	}
}

func TestGeneraService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Genus{{Id: 1, GenusName: "Test Genus"}}

	var called bool
	mux.HandleFunc(urlPath(t, router.Genera, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	genera, err := client.Genera.List(nil)
	if err != nil {
		t.Errorf("Genera.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}
	if !reflect.DeepEqual(genera, want) {
		t.Errorf("Genera.List return %+v, want %+v", genera, want)
	}
}

func TestGeneraService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := &Genus{Id: 1, GenusName: "Test Genus"}

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateGenus, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"genus_name":"Test Genus Updated","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	genus := &Genus{Id: 1, GenusName: "Test Genus Updated"}
	updated, err := client.Genera.Update(1, genus)
	if err != nil {
		t.Errorf("Genera.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

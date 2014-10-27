package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newSpecies() *Species {
	species := NewSpecies()
	species.Id = 1
	species.GenusId = 1
	return species
}

func TestSpeciesService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newSpecies()

	var called bool
	mux.HandleFunc(urlPath(t, router.Species, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	species, err := client.Species.Get(want.Id)
	if err != nil {
		t.Errorf("Species.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(species, want) {
		t.Errorf("Species.Get returned %+v, want %+v", species, want)
	}
}

func TestSpeciesService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newSpecies()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateSpecies, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"genus_id":1,"species_name":"Test Species","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	species := newSpecies()
	created, err := client.Species.Create(species)
	if err != nil {
		t.Errorf("Species.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(species, want) {
		t.Errorf("Species.Create returned %+v, want %+v", species, want)
	}
}

func TestSpeciesService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Species{newSpecies()}

	var called bool
	mux.HandleFunc(urlPath(t, router.SpeciesList, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	species, err := client.Species.List(nil)
	if err != nil {
		t.Errorf("Species.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(species, want) {
		t.Errorf("Species.List return %+v, want %+v", species, want)
	}
}

func TestSpeciesService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newSpecies()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateSpecies, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"genus_id":1,"species_name":"Test Species Updated","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	species := newSpecies()
	species.SpeciesName = "Test Species Updated"
	updated, err := client.Species.Update(species.Id, species)
	if err != nil {
		t.Errorf("Species.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestSpeciesService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newSpecies()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteSpecies, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.Species.Delete(want.Id)
	if err != nil {
		t.Errorf("Species.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

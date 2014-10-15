package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func TestSpeciesService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := &Species{Id: 1, GenusId: 1, SpeciesName: "Test Species"}

	var called bool
	mux.HandleFunc(urlPath(t, router.Species, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	species, err := client.Species.Get(1)
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

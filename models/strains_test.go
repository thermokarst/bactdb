package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newStrain() *Strain {
	strain := NewStrain()
	strain.Id = 1
	strain.SpeciesId = 1
	return strain
}

func TestStrainService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newStrain()

	var called bool
	mux.HandleFunc(urlPath(t, router.Strain, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	strain, err := client.Strains.Get(want.Id)
	if err != nil {
		t.Errorf("Strain.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(strain, want) {
		t.Errorf("Strain.Get return %+v, want %+v", strain, want)
	}
}

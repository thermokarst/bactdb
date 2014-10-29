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

func TestStrainService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newStrain()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateStrain, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"species_id":1,"strain_name":"Test Strain","strain_type":"Test Type","etymology":"Test Etymology","accession_banks":"Test Accession","genbank_eml_ddb":"Test Genbank","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z","deleted_at":"0001-01-01T00:00:00Z"}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	strain := newStrain()
	created, err := client.Strains.Create(strain)
	if err != nil {
		t.Errorf("Strains.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(strain, want) {
		t.Errorf("Strains.Create returned %+v, want %+v", strain, want)
	}
}

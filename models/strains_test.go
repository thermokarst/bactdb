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
		testBody(t, r, `{"id":1,"speciesId":1,"strainName":"Test Strain","strainType":"Test Type","etymology":"Test Etymology","accessionBanks":"Test Accession","genbankEmblDdb":"Test Genbank","isolatedFrom":null,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")

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

func TestStrainService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Strain{newStrain()}

	var called bool
	mux.HandleFunc(urlPath(t, router.Strains, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	strains, err := client.Strains.List(nil)
	if err != nil {
		t.Errorf("Strains.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(strains, want) {
		t.Errorf("Strains.List return %+v, want %+v", strains, want)
	}
}

func TestStrainService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newStrain()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateStrain, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"speciesId":1,"strainName":"Test Strain Updated","strainType":"Test Type Updated","etymology":"Test Etymology Updated","accessionBanks":"Test Accession Updated","genbankEmblDdb":"Test Genbank Updated","isolatedFrom":null,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	strain := newStrain()
	strain.StrainName = "Test Strain Updated"
	strain.StrainType = "Test Type Updated"
	strain.Etymology = "Test Etymology Updated"
	strain.AccessionBanks = "Test Accession Updated"
	strain.GenbankEmblDdb = "Test Genbank Updated"
	updated, err := client.Strains.Update(strain.Id, strain)
	if err != nil {
		t.Errorf("Strains.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestStrainService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newStrain()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteStrain, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.Strains.Delete(want.Id)
	if err != nil {
		t.Errorf("Strains.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

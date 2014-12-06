package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newUnitType() *UnitType {
	unit_type := NewUnitType()
	unit_type.Id = 1
	return unit_type
}

func TestUnitTypeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newUnitType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UnitType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	unit_type, err := client.UnitTypes.Get(want.Id)
	if err != nil {
		t.Errorf("UnitTypes.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(unit_type, want) {
		t.Errorf("UnitTypes.Get return %+v, want %+v", unit_type, want)
	}
}

func TestUnitTypeService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newUnitType()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateUnitType, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"name":"Test Unit Type","symbol":"x","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	unit_type := newUnitType()
	created, err := client.UnitTypes.Create(unit_type)
	if err != nil {
		t.Errorf("UnitTypes.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(unit_type, want) {
		t.Errorf("UnitTypes.Create returned %+v, want %+v", unit_type, want)
	}
}

func TestUnitTypeService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*UnitType{newUnitType()}

	var called bool
	mux.HandleFunc(urlPath(t, router.UnitTypes, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	unit_types, err := client.UnitTypes.List(nil)
	if err != nil {
		t.Errorf("UnitTypes.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(unit_types, want) {
		t.Errorf("UnitTypes.List return %+v, want %+v", unit_types, want)
	}
}

func TestUnitTypeService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newUnitType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateUnitType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"name":"Test Unit Type Updated","symbol":"x","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	unit_type := newUnitType()
	unit_type.Name = "Test Unit Type Updated"
	updated, err := client.UnitTypes.Update(unit_type.Id, unit_type)
	if err != nil {
		t.Errorf("UnitTypes.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestUnitTypeService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newUnitType()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteUnitType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.UnitTypes.Delete(want.Id)
	if err != nil {
		t.Errorf("UnitTypes.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

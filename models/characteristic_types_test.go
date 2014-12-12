package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newCharacteristicType() *CharacteristicType {
	characteristic_type := NewCharacteristicType()
	characteristic_type.Id = 1
	return characteristic_type
}

func TestCharacteristicTypeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristicType()

	var called bool
	mux.HandleFunc(urlPath(t, router.CharacteristicType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	characteristic_type, err := client.CharacteristicTypes.Get(want.Id)
	if err != nil {
		t.Errorf("CharacteristicTypes.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(characteristic_type, want) {
		t.Errorf("CharacteristicTypes.Get return %+v, want %+v", characteristic_type, want)
	}
}

func TestCharacteristicTypeService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristicType()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateCharacteristicType, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"characteristicTypeName":"Test Obs Type","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	characteristic_type := newCharacteristicType()
	created, err := client.CharacteristicTypes.Create(characteristic_type)
	if err != nil {
		t.Errorf("CharacteristicTypes.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(characteristic_type, want) {
		t.Errorf("CharacteristicTypes.Create returned %+v, want %+v", characteristic_type, want)
	}
}

func TestCharacteristicTypeService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*CharacteristicType{newCharacteristicType()}

	var called bool
	mux.HandleFunc(urlPath(t, router.CharacteristicTypes, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	characteristic_types, err := client.CharacteristicTypes.List(nil)
	if err != nil {
		t.Errorf("CharacteristicTypes.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(characteristic_types, want) {
		t.Errorf("CharacteristicTypes.List return %+v, want %+v", characteristic_types, want)
	}
}

func TestCharacteristicTypeService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristicType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateCharacteristicType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"characteristicTypeName":"Test Obs Type Updated","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	characteristic_type := newCharacteristicType()
	characteristic_type.CharacteristicTypeName = "Test Obs Type Updated"
	updated, err := client.CharacteristicTypes.Update(characteristic_type.Id, characteristic_type)
	if err != nil {
		t.Errorf("CharacteristicTypes.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestCharacteristicTypeService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristicType()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteCharacteristicType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.CharacteristicTypes.Delete(want.Id)
	if err != nil {
		t.Errorf("CharacteristicTypes.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

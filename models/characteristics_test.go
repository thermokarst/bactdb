package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newCharacteristic() *Characteristic {
	characteristic := NewCharacteristic()
	characteristic.Id = 1
	return characteristic
}

func TestCharacteristicService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristic()

	var called bool
	mux.HandleFunc(urlPath(t, router.Characteristic, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, CharacteristicJSON{Characteristic: want})
	})

	characteristic, err := client.Characteristics.Get(want.Id)
	if err != nil {
		t.Errorf("Characteristics.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(characteristic, want) {
		t.Errorf("Characteristics.Get return %+v, want %+v", characteristic, want)
	}
}

func TestCharacteristicService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristic()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateCharacteristic, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"characteristic":{"id":1,"characteristicName":"Test Characteristic","characteristicTypeId":0,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null,"measurements":[]}}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	characteristic := newCharacteristic()
	created, err := client.Characteristics.Create(characteristic)
	if err != nil {
		t.Errorf("Characteristics.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(characteristic, want) {
		t.Errorf("Characteristics.Create returned %+v, want %+v", characteristic, want)
	}
}

func TestCharacteristicService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Characteristic{newCharacteristic()}

	var called bool
	mux.HandleFunc(urlPath(t, router.Characteristics, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, CharacteristicsJSON{Characteristics: want})
	})

	characteristics, err := client.Characteristics.List(nil)
	if err != nil {
		t.Errorf("Characteristics.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(characteristics, want) {
		t.Errorf("Characteristics.List return %+v, want %+v", characteristics, want)
	}
}

func TestCharacteristicService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristic()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateCharacteristic, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"characteristic":{"id":1,"characteristicName":"Test Char Updated","characteristicTypeId":0,"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":null,"measurements":[]}}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	characteristic := newCharacteristic()
	characteristic.CharacteristicName = "Test Char Updated"
	updated, err := client.Characteristics.Update(characteristic.Id, characteristic)
	if err != nil {
		t.Errorf("Characteristics.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestCharacteristicService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newCharacteristic()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteCharacteristic, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")
		testBody(t, r, "")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.Characteristics.Delete(want.Id)
	if err != nil {
		t.Errorf("Characteristics.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

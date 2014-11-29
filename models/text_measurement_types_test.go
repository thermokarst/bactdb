package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newTextMeasurementType() *TextMeasurementType {
	text_measurement_type := NewTextMeasurementType()
	text_measurement_type.Id = 1
	return text_measurement_type
}

func TestTextMeasurementTypeService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newTextMeasurementType()

	var called bool
	mux.HandleFunc(urlPath(t, router.TextMeasurementType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	text_measurement_type, err := client.TextMeasurementTypes.Get(want.Id)
	if err != nil {
		t.Errorf("TextMeasurementTypes.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(text_measurement_type, want) {
		t.Errorf("TextMeasurementTypes.Get return %+v, want %+v", text_measurement_type, want)
	}
}

func TestTextMeasurementTypeService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newTextMeasurementType()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateTextMeasurementType, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"textMeasurementName":"Test Text Measurement Type","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	text_measurement_type := newTextMeasurementType()
	created, err := client.TextMeasurementTypes.Create(text_measurement_type)
	if err != nil {
		t.Errorf("TextMeasurementTypes.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(text_measurement_type, want) {
		t.Errorf("TextMeasurementTypes.Create returned %+v, want %+v", text_measurement_type, want)
	}
}

func TestTextMeasurementTypeService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*TextMeasurementType{newTextMeasurementType()}

	var called bool
	mux.HandleFunc(urlPath(t, router.TextMeasurementTypes, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	text_measurement_type, err := client.TextMeasurementTypes.List(nil)
	if err != nil {
		t.Errorf("TextMeasurementTypes.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(text_measurement_type, want) {
		t.Errorf("TextMeasurementTypes.List return %+v, want %+v", text_measurement_type, want)
	}
}

func TestTextMeasurementTypeService_Update(t *testing.T) {
	setup()
	defer teardown()

	want := newTextMeasurementType()

	var called bool
	mux.HandleFunc(urlPath(t, router.UpdateTextMeasurementType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "PUT")
		testBody(t, r, `{"id":1,"textMeasurementName":"Test Text Measurement Type Updated","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")
		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	text_measurement_type := newTextMeasurementType()
	text_measurement_type.TextMeasurementName = "Test Text Measurement Type Updated"
	updated, err := client.TextMeasurementTypes.Update(text_measurement_type.Id, text_measurement_type)
	if err != nil {
		t.Errorf("TextMeasurementTypes.Update returned error: %v", err)
	}

	if !updated {
		t.Error("!updated")
	}

	if !called {
		t.Fatal("!called")
	}
}

func TestTextMeasurementTypeService_Delete(t *testing.T) {
	setup()
	defer teardown()

	want := newTextMeasurementType()

	var called bool
	mux.HandleFunc(urlPath(t, router.DeleteTextMeasurementType, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "DELETE")

		w.WriteHeader(http.StatusOK)
		writeJSON(w, want)
	})

	deleted, err := client.TextMeasurementTypes.Delete(want.Id)
	if err != nil {
		t.Errorf("TextMeasurementTypes.Delete returned error: %v", err)
	}

	if !deleted {
		t.Error("!deleted")
	}

	if !called {
		t.Fatal("!called")
	}
}

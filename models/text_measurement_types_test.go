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

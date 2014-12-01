package models

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newMeasurement() *Measurement {
	measurement := NewMeasurement()
	measurement.Id = 1
	return measurement
}

func TestMeasurementService_Get(t *testing.T) {
	setup()
	defer teardown()

	want := newMeasurement()

	var called bool
	mux.HandleFunc(urlPath(t, router.Measurement, map[string]string{"Id": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	measurement, err := client.Measurements.Get(want.Id)
	if err != nil {
		t.Errorf("Measurements.Get returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(measurement, want) {
		t.Errorf("Measurements.Get return %+v, want %+v", measurement, want)
	}
}

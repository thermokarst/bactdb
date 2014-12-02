package models

import (
	"database/sql"
	"net/http"
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/router"
)

func newMeasurement() *Measurement {
	measurement := NewMeasurement()
	measurement.Id = 1
	measurement.StrainId = 1
	measurement.ObservationId = 1
	measurement.UnitTypeId = sql.NullInt64{Int64: 1, Valid: true}
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

func TestMeasurementService_Create(t *testing.T) {
	setup()
	defer teardown()

	want := newMeasurement()

	var called bool
	mux.HandleFunc(urlPath(t, router.CreateMeasurement, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "POST")
		testBody(t, r, `{"id":1,"strainId":1,"observationId":1,"textMeasurementTypeId":{"Int64":0,"Valid":false},"measurementValue":{"Float64":1.23,"Valid":true},"confidenceInterval":{"Float64":0,"Valid":false},"unitTypeId":{"Int64":1,"Valid":true},"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","deletedAt":{"Time":"0001-01-01T00:00:00Z","Valid":false}}`+"\n")

		w.WriteHeader(http.StatusCreated)
		writeJSON(w, want)
	})

	measurement := newMeasurement()
	created, err := client.Measurements.Create(measurement)
	if err != nil {
		t.Errorf("Measurements.Create returned error: %v", err)
	}

	if !created {
		t.Error("!created")
	}

	if !called {
		t.Fatal("!called")
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(measurement, want) {
		t.Errorf("Measurements.Create returned %+v, want %+v", measurement, want)
	}
}

func TestMeasurementService_List(t *testing.T) {
	setup()
	defer teardown()

	want := []*Measurement{newMeasurement()}

	var called bool
	mux.HandleFunc(urlPath(t, router.Measurements, nil), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{})

		writeJSON(w, want)
	})

	measurements, err := client.Measurements.List(nil)
	if err != nil {
		t.Errorf("Measurements.List returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	for _, u := range want {
		normalizeTime(&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	}

	if !reflect.DeepEqual(measurements, want) {
		t.Errorf("Measurements.List return %+v, want %+v", measurements, want)
	}
}
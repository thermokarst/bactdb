package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newMeasurement() *models.Measurement {
	measurement := models.NewMeasurement()
	return measurement
}

func TestMeasurement_Get(t *testing.T) {
	setup()

	want := newMeasurement()

	calledGet := false

	store.Measurements.(*models.MockMeasurementsService).Get_ = func(id int64) (*models.Measurement, error) {
		if id != want.Id {
			t.Errorf("wanted request for measurement %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Measurements.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got %+v but wanted %+v", got, want)
	}
}

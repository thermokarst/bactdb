package api

import (
	"database/sql"
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newMeasurement() *models.Measurement {
	measurement := models.NewMeasurement()
	measurement.Id = 1
	measurement.StrainId = 2
	measurement.CharacteristicId = 3
	measurement.TextMeasurementTypeId = models.NullInt64{sql.NullInt64{Int64: 4, Valid: false}}
	measurement.UnitTypeId = models.NullInt64{sql.NullInt64{Int64: 5, Valid: true}}
	measurement.Notes = models.NullString{sql.NullString{String: "a note", Valid: true}}
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

func TestMeasurement_Create(t *testing.T) {
	setup()

	want := newMeasurement()

	calledPost := false
	store.Measurements.(*models.MockMeasurementsService).Create_ = func(measurement *models.Measurement) (bool, error) {
		if !normalizeDeepEqual(want, measurement) {
			t.Errorf("wanted request for measurement %d but got %d", want, measurement)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Measurements.Create(want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPost {
		t.Error("!calledPost")
	}
	if !success {
		t.Error("!success")
	}
}

func TestMeasurement_List(t *testing.T) {
	setup()

	want := []*models.Measurement{newMeasurement()}
	wantOpt := &models.MeasurementListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Measurements.(*models.MockMeasurementsService).List_ = func(opt *models.MeasurementListOptions) ([]*models.Measurement, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	measurements, err := apiClient.Measurements.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &measurements) {
		t.Errorf("got measurements %+v but wanted measurements %+v", measurements, want)
	}
}

func TestMeasurement_Update(t *testing.T) {
	setup()

	want := newMeasurement()

	calledPut := false
	store.Measurements.(*models.MockMeasurementsService).Update_ = func(id int64, measurement *models.Measurement) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for measurement %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, measurement) {
			t.Errorf("wanted request for measurement %d but got %d", want, measurement)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.Measurements.Update(want.Id, want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPut {
		t.Error("!calledPut")
	}
	if !success {
		t.Error("!success")
	}
}

func TestMeasurement_Delete(t *testing.T) {
	setup()

	want := newMeasurement()

	calledDelete := false
	store.Measurements.(*models.MockMeasurementsService).Delete_ = func(id int64) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for measurement %d but got %d", want.Id, id)
		}
		calledDelete = true
		return true, nil
	}

	success, err := apiClient.Measurements.Delete(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledDelete {
		t.Error("!calledDelete")
	}
	if !success {
		t.Error("!success")
	}
}

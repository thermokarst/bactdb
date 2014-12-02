package datastore

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertMeasurement(t *testing.T, tx *modl.Transaction) *models.Measurement {
	// clean up our target table
	tx.Exec(`DELETE FROM measurements;`)
	measurement := newMeasurement(t, tx)
	if err := tx.Insert(measurement); err != nil {
		t.Fatal(err)
	}
	return measurement
}

func newMeasurement(t *testing.T, tx *modl.Transaction) *models.Measurement {
	// we have a few things to take care of first...
	strain := insertStrain(t, tx)
	observation := insertObservation(t, tx)

	// we want to create and insert a unit type record, too.
	unit_type := insertUnitType(t, tx)

	return &models.Measurement{
		StrainId:         strain.Id,
		ObservationId:    observation.Id,
		MeasurementValue: sql.NullFloat64{Float64: 1.23, Valid: true},
		UnitTypeId:       sql.NullInt64{Int64: unit_type.Id, Valid: true},
	}
}

func TestMeasurementsStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertMeasurement(t, tx)

	d := NewDatastore(tx)

	measurement, err := d.Measurements.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&measurement.CreatedAt, &measurement.UpdatedAt, &measurement.DeletedAt)

	if !reflect.DeepEqual(measurement, want) {
		t.Errorf("got measurement %+v, want %+v", measurement, want)
	}
}

func TestMeasurementsStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	measurement := newMeasurement(t, tx)

	d := NewDatastore(tx)

	created, err := d.Measurements.Create(measurement)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if measurement.Id == 0 {
		t.Error("want nonzero measurement.Id after submitting")
	}
}

func TestMeasurementsStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_measurement := insertMeasurement(t, tx)
	want := []*models.Measurement{want_measurement}

	d := NewDatastore(tx)

	measurements, err := d.Measurements.List(&models.MeasurementListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&measurements[i].CreatedAt, &measurements[i].UpdatedAt, &measurements[i].DeletedAt)
	}
	if !reflect.DeepEqual(measurements, want) {
		t.Errorf("got measurements %+v, want %+v", measurements, want)
	}
}

func TestMeasurementsStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	measurement := insertMeasurement(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	measurement.MeasurementValue = sql.NullFloat64{Float64: 4.56, Valid: true}
	updated, err := d.Measurements.Update(measurement.Id, measurement)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestMeasurementsStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	measurement := insertMeasurement(t, tx)

	d := NewDatastore(tx)

	// Delete it
	deleted, err := d.Measurements.Delete(measurement.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

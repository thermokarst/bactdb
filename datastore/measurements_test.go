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

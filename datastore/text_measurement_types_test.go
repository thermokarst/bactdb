package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertTextMeasurementType(t *testing.T, tx *modl.Transaction) *models.TextMeasurementType {
	// clean up our target table
	tx.Exec(`DELETE FROM text_measurement_types;`)
	text_measurement_type := newTextMeasurementType(t, tx)
	if err := tx.Insert(text_measurement_type); err != nil {
		t.Fatal(err)
	}
	return text_measurement_type
}

func newTextMeasurementType(t *testing.T, tx *modl.Transaction) *models.TextMeasurementType {
	return &models.TextMeasurementType{TextMeasurementName: "Test Text Measurement Type"}
}

func TestTextMeasurementTypesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertTextMeasurementType(t, tx)

	d := NewDatastore(tx)

	text_measurement_type, err := d.TextMeasurementTypes.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&text_measurement_type.CreatedAt, &text_measurement_type.UpdatedAt, &text_measurement_type.DeletedAt)

	if !reflect.DeepEqual(text_measurement_type, want) {
		t.Errorf("got text_measurement_type %+v, want %+v", text_measurement_type, want)
	}
}

func TestTextMeasurementTypesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	text_measurement_type := newTextMeasurementType(t, tx)

	d := NewDatastore(tx)

	created, err := d.TextMeasurementTypes.Create(text_measurement_type)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if text_measurement_type.Id == 0 {
		t.Error("want nonzero text_measurement_type.Id after submitting")
	}
}

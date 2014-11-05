package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertObservationType(t *testing.T, tx *modl.Transaction) *models.ObservationType {
	// clean up our target table
	tx.Exec(`DELETE FROM observation_types;`)
	observation_type := newObservationType(t, tx)
	if err := tx.Insert(observation_type); err != nil {
		t.Fatal(err)
	}
	return observation_type
}

func newObservationType(t *testing.T, tx *modl.Transaction) *models.ObservationType {
	return &models.ObservationType{ObservationTypeName: "Test Obs"}
}

func TestObservationTypesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertObservationType(t, tx)

	d := NewDatastore(tx)

	observation_type, err := d.ObservationTypes.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&observation_type.CreatedAt, &observation_type.UpdatedAt, &observation_type.DeletedAt)

	if !reflect.DeepEqual(observation_type, want) {
		t.Errorf("got observation_type %+v, want %+v", observation_type, want)
	}
}

func TestObservationTypesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	observation_type := newObservationType(t, tx)

	d := NewDatastore(tx)

	created, err := d.ObservationTypes.Create(observation_type)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if observation_type.Id == 0 {
		t.Error("want nonzero observation_type.Id after submitting")
	}
}

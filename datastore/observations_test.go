package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertObservation(t *testing.T, tx *modl.Transaction) *models.Observation {
	// clean up our target table
	tx.Exec(`DELETE FROM observations;`)
	observation := newObservation(t, tx)
	if err := tx.Insert(observation); err != nil {
		t.Fatal(err)
	}
	return observation
}

func newObservation(t *testing.T, tx *modl.Transaction) *models.Observation {
	// we want to create and insert an observation type record, too.
	observation_type := insertObservationType(t, tx)
	return &models.Observation{ObservationName: "Test Observation",
		ObservationTypeId: observation_type.Id}
}

func TestObservationsStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertObservation(t, tx)

	d := NewDatastore(tx)

	observation, err := d.Observations.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&observation.CreatedAt, &observation.UpdatedAt, &observation.DeletedAt)

	if !reflect.DeepEqual(observation, want) {
		t.Errorf("got observation %+v, want %+v", observation, want)
	}
}

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

func TestObservationsStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	observation := newObservation(t, tx)

	d := NewDatastore(tx)

	created, err := d.Observations.Create(observation)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if observation.Id == 0 {
		t.Error("want nonzero observation.Id after submitting")
	}
}

func TestObservationsStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_observation := insertObservation(t, tx)
	want := []*models.Observation{want_observation}

	d := NewDatastore(tx)

	observations, err := d.Observations.List(&models.ObservationListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&observations[i].CreatedAt, &observations[i].UpdatedAt, &observations[i].DeletedAt)
	}
	if !reflect.DeepEqual(observations, want) {
		t.Errorf("got observations %+v, want %+v", observations, want)
	}
}

func TestObservationsStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	observation := insertObservation(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	observation.ObservationName = "Updated Obs"
	updated, err := d.Observations.Update(observation.Id, observation)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

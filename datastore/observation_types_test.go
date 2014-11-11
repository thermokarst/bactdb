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

func TestObservationTypesStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_observation_type := insertObservationType(t, tx)
	want := []*models.ObservationType{want_observation_type}

	d := NewDatastore(tx)

	observation_types, err := d.ObservationTypes.List(&models.ObservationTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&observation_types[i].CreatedAt, &observation_types[i].UpdatedAt, &observation_types[i].DeletedAt)
	}
	if !reflect.DeepEqual(observation_types, want) {
		t.Errorf("got observation_types %+v, want %+v", observation_types, want)
	}
}

func TestObservationTypesStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	observation_type := insertObservationType(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	observation_type.ObservationTypeName = "Updated Obs Type"
	updated, err := d.ObservationTypes.Update(observation_type.Id, observation_type)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

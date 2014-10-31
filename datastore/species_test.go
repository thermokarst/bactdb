package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertSpecies(t *testing.T, tx *modl.Transaction) *models.Species {
	// clean up our target table
	tx.Exec(`DELETE FROM species;`)
	species := newSpecies(t, tx)
	if err := tx.Insert(species); err != nil {
		t.Fatal(err)
	}
	return species
}

func newSpecies(t *testing.T, tx *modl.Transaction) *models.Species {
	// we want to create and insert a genus record, too
	genus := insertGenus(t, tx)
	return &models.Species{GenusId: genus.Id, SpeciesName: "Test Species"}
}

func TestSpeciesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertSpecies(t, tx)

	d := NewDatastore(tx)

	species, err := d.Species.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&species.CreatedAt, &species.UpdatedAt, &species.DeletedAt)
	if !reflect.DeepEqual(species, want) {
		t.Errorf("got species %+v, want %+v", species, want)
	}
}

func TestSpeciesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	species := newSpecies(t, tx)

	d := NewDatastore(tx)

	created, err := d.Species.Create(species)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if species.Id == 0 {
		t.Error("want nonzero species.Id after submitting")
	}
}

func TestSpeciesStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_species := insertSpecies(t, tx)
	want := []*models.Species{want_species}

	d := NewDatastore(tx)

	species, err := d.Species.List(&models.SpeciesListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&species[i].CreatedAt, &species[i].UpdatedAt, &species[i].DeletedAt)
	}
	if !reflect.DeepEqual(species, want) {
		t.Errorf("got species %+v, want %+v", species, want)
	}
}

func TestSpeciesStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	species := insertSpecies(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	species.SpeciesName = "Updated Species"
	updated, err := d.Species.Update(species.Id, species)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestSpeciesStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	species := insertSpecies(t, tx)

	d := NewDatastore(tx)

	// Delete it
	deleted, err := d.Species.Delete(species.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

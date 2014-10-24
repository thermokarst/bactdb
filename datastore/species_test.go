package datastore

import (
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func TestSpeciesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM species;`)
	tx.Exec(`DELETE FROM genera;`)

	wantGenus := &models.Genus{GenusName: "Test Genus"}
	if err := tx.Insert(wantGenus); err != nil {
		t.Fatal(err)
	}

	want := &models.Species{Id: 1, GenusId: wantGenus.Id, SpeciesName: "Test Species"}
	if err := tx.Insert(want); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	species, err := d.Species.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(species, want) {
		t.Errorf("got species %+v, want %+v", species, want)
	}
}

func TestSpeciesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM species;`)
	tx.Exec(`DELETE FROM genera;`)

	genus := &models.Genus{}
	if err := tx.Insert(genus); err != nil {
		t.Fatal(err)
	}

	species := &models.Species{Id: 1, GenusId: genus.Id, SpeciesName: "Test Species"}

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

	// Test on a clean database
	tx.Exec(`DELETE FROM species;`)
	tx.Exec(`DELETE FROM genera;`)

	genus := &models.Genus{}

	if err := tx.Insert(genus); err != nil {
		t.Fatal(err)
	}

	want := []*models.Species{{Id: 1, GenusId: genus.Id, SpeciesName: "Test Species"}}

	if err := tx.Insert(want[0]); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	species, err := d.Species.List(&models.SpeciesListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for _, g := range want {
		normalizeTime(&g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
	}
	if !reflect.DeepEqual(species, want) {
		t.Errorf("got species %+v, want %+v", species, want)
	}
}

func TestSpeciesStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM species;`)
	tx.Exec(`DELETE FROM genera;`)

	d := NewDatastore(nil)
	// Add a new record
	genus := &models.Genus{GenusName: "Test Genus"}
	_, err := d.Genera.Create(genus)
	if err != nil {
		t.Fatal(err)
	}
	species := &models.Species{GenusId: genus.Id, SpeciesName: "Test Species"}
	created, err := d.Species.Create(species)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}

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

	// Test on a clean database
	tx.Exec(`DELETE FROM species;`)
	tx.Exec(`DELETE FROM genera;`)

	d := NewDatastore(tx)
	// Add a new record
	genus := &models.Genus{GenusName: "Test Genus"}
	_, err := d.Genera.Create(genus)
	if err != nil {
		t.Fatal(err)
	}
	species := &models.Species{GenusId: genus.Id, SpeciesName: "Test Species"}
	created, err := d.Species.Create(species)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}

	// Delete it
	deleted, err := d.Species.Delete(species.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

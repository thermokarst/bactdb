package datastore

import (
	"reflect"
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func TestGeneraStore_Get_db(t *testing.T) {
	want := &models.Genus{Id: 1, GenusName: "Test Genus"}

	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM genera;`)
	if err := tx.Insert(want); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	genus, err := d.Genera.Get(1)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	if !reflect.DeepEqual(genus, want) {
		t.Errorf("got genus %+v, want %+v", genus, want)
	}
}

func TestGeneraStore_Create_db(t *testing.T) {
	genus := &models.Genus{Id: 1, GenusName: "Test Genus"}

	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM genera;`)

	d := NewDatastore(tx)
	created, err := d.Genera.Create(genus)
	if err != nil {
		t.Fatal(err)
	}

	if !created {
		t.Error("!created")
	}
	if genus.Id == 0 {
		t.Error("want nonzero genus.Id after submitting")
	}
}

func TestGeneraStore_List_db(t *testing.T) {
	want := []*models.Genus{{Id: 1, GenusName: "Test Genus"}}

	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM genera;`)
	if err := tx.Insert(want[0]); err != nil {
		t.Fatal(err)
	}

	d := NewDatastore(tx)
	genera, err := d.Genera.List(&models.GenusListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for _, g := range want {
		normalizeTime(&g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
	}
	if !reflect.DeepEqual(genera, want) {
		t.Errorf("got genera %+v, want %+v", genera, want)
	}
}

func TestGeneraStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM genera;`)

	d := NewDatastore(nil)
	// Add a new record
	genus := &models.Genus{GenusName: "Test Genus"}
	created, err := d.Genera.Create(genus)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}

	// Tweak it
	genus.GenusName = "Updated Genus"
	updated, err := d.Genera.Update(genus.Id, genus)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestGeneraStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	// Test on a clean database
	tx.Exec(`DELETE FROM genera;`)

	d := NewDatastore(tx)
	// Add a new record
	genus := &models.Genus{GenusName: "Test Genus"}
	created, err := d.Genera.Create(genus)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}

	// Delete it
	deleted, err := d.Genera.Delete(genus.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

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

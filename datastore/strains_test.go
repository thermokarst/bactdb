package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertStrain(t *testing.T, tx *modl.Transaction) *models.Strain {
	// clean up our target table
	tx.Exec(`DELETE FROM strains;`)
	strain := newStrain(t, tx)
	if err := tx.Insert(strain); err != nil {
		t.Fatal(err)
	}
	return strain
}

func newStrain(t *testing.T, tx *modl.Transaction) *models.Strain {
	// we want to create and insert a species (and genus) record too
	species := insertSpecies(t, tx)
	return &models.Strain{SpeciesId: species.Id, StrainName: "Test Strain",
		StrainType: "Test Type", Etymology: "Test Etymology",
		AccessionBanks: "Test Bank", GenbankEmblDdb: "Test Genbank"}
}

func TestStrainsStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertStrain(t, tx)

	d := NewDatastore(tx)

	strain, err := d.Strains.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)

	if !reflect.DeepEqual(strain, want) {
		t.Errorf("got strain %+v, want %+v", strain, want)
	}
}

func TestStrainsStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	strain := newStrain(t, tx)

	d := NewDatastore(tx)

	created, err := d.Strains.Create(strain)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if strain.Id == 0 {
		t.Error("want nonzero strain.Id after submitting")
	}
}

func TestStrainsStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_strain := insertStrain(t, tx)
	want := []*models.Strain{want_strain}

	d := NewDatastore(tx)

	strains, err := d.Strains.List(&models.StrainListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for _, g := range want {
		normalizeTime(&g.CreatedAt, &g.UpdatedAt, &g.DeletedAt)
	}
	if !reflect.DeepEqual(strains, want) {
		t.Errorf("got strains %+v, want %+v", strains, want)
	}
}

func TestStrainsStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	strain := insertStrain(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	strain.StrainName = "Updated Strain"
	updated, err := d.Strains.Update(strain.Id, strain)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestStrainsStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	strain := insertStrain(t, tx)

	d := NewDatastore(tx)

	// Delete it
	deleted, err := d.Strains.Delete(strain.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

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

package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertUnitType(t *testing.T, tx *modl.Transaction) *models.UnitType {
	// clean up our target table
	tx.Exec(`DELETE FROM unit_types;`)
	unit_type := newUnitType(t, tx)
	if err := tx.Insert(unit_type); err != nil {
		t.Fatal(err)
	}
	return unit_type
}

func newUnitType(t *testing.T, tx *modl.Transaction) *models.UnitType {
	return &models.UnitType{Name: "Test Unit Type", Symbol: "x"}
}

func TestUnitTypesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertUnitType(t, tx)

	d := NewDatastore(tx)

	unit_type, err := d.UnitTypes.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&unit_type.CreatedAt, &unit_type.UpdatedAt, &unit_type.DeletedAt)

	if !reflect.DeepEqual(unit_type, want) {
		t.Errorf("got unit_type %+v, want %+v", unit_type, want)
	}
}

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

func TestUnitTypesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	unit_type := newUnitType(t, tx)

	d := NewDatastore(tx)

	created, err := d.UnitTypes.Create(unit_type)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if unit_type.Id == 0 {
		t.Error("want nonzero unit_type.Id after submitting")
	}
}

func TestUnitTypesStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_unit_type := insertUnitType(t, tx)
	want := []*models.UnitType{want_unit_type}

	d := NewDatastore(tx)

	unit_types, err := d.UnitTypes.List(&models.UnitTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&unit_types[i].CreatedAt, &unit_types[i].UpdatedAt, &unit_types[i].DeletedAt)
	}
	if !reflect.DeepEqual(unit_types, want) {
		t.Errorf("got unit_types %+v, want %+v", unit_types, want)
	}
}

func TestUnitTypesStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	unit_type := insertUnitType(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	unit_type.Name = "Updated Unit Type"
	updated, err := d.UnitTypes.Update(unit_type.Id, unit_type)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

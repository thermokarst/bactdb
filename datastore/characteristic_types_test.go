package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertCharacteristicType(t *testing.T, tx *modl.Transaction) *models.CharacteristicType {
	// clean up our target table
	tx.Exec(`DELETE FROM characteristic_types;`)
	characteristic_type := newCharacteristicType(t, tx)
	if err := tx.Insert(characteristic_type); err != nil {
		t.Fatal(err)
	}
	return characteristic_type
}

func newCharacteristicType(t *testing.T, tx *modl.Transaction) *models.CharacteristicType {
	return &models.CharacteristicType{CharacteristicTypeName: "Test Obs"}
}

func TestCharacteristicTypesStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertCharacteristicType(t, tx)

	d := NewDatastore(tx)

	characteristic_type, err := d.CharacteristicTypes.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&characteristic_type.CreatedAt, &characteristic_type.UpdatedAt, &characteristic_type.DeletedAt)

	if !reflect.DeepEqual(characteristic_type, want) {
		t.Errorf("got characteristic_type %+v, want %+v", characteristic_type, want)
	}
}

func TestCharacteristicTypesStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	characteristic_type := newCharacteristicType(t, tx)

	d := NewDatastore(tx)

	created, err := d.CharacteristicTypes.Create(characteristic_type)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if characteristic_type.Id == 0 {
		t.Error("want nonzero characteristic_type.Id after submitting")
	}
}

func TestCharacteristicTypesStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_characteristic_type := insertCharacteristicType(t, tx)
	want := []*models.CharacteristicType{want_characteristic_type}

	d := NewDatastore(tx)

	characteristic_types, err := d.CharacteristicTypes.List(&models.CharacteristicTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&characteristic_types[i].CreatedAt, &characteristic_types[i].UpdatedAt, &characteristic_types[i].DeletedAt)
	}
	if !reflect.DeepEqual(characteristic_types, want) {
		t.Errorf("got characteristic_types %+v, want %+v", characteristic_types, want)
	}
}

func TestCharacteristicTypesStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	characteristic_type := insertCharacteristicType(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	characteristic_type.CharacteristicTypeName = "Updated Obs Type"
	updated, err := d.CharacteristicTypes.Update(characteristic_type.Id, characteristic_type)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestCharacteristicTypesStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	characteristic_type := insertCharacteristicType(t, tx)

	d := NewDatastore(tx)

	// Delete it
	deleted, err := d.CharacteristicTypes.Delete(characteristic_type.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

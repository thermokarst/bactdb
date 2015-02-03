package datastore

import (
	"reflect"
	"testing"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

func insertCharacteristic(t *testing.T, tx *modl.Transaction) *models.Characteristic {
	// clean up our target table
	tx.Exec(`DELETE FROM characteristics;`)
	c := newCharacteristic(t, tx)
	if err := tx.Insert(c); err != nil {
		t.Fatal(err)
	}
	return &models.Characteristic{c, []int64(nil)}
}

func newCharacteristic(t *testing.T, tx *modl.Transaction) *models.CharacteristicBase {
	// we want to create and insert an characteristic type record, too.
	characteristic_type := insertCharacteristicType(t, tx)
	return &models.CharacteristicBase{CharacteristicName: "Test Characteristic",
		CharacteristicTypeId: characteristic_type.Id}
}

func TestCharacteristicsStore_Get_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want := insertCharacteristic(t, tx)

	d := NewDatastore(tx)

	characteristic, err := d.Characteristics.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	normalizeTime(&want.CreatedAt, &want.UpdatedAt, &want.DeletedAt)
	normalizeTime(&characteristic.CreatedAt, &characteristic.UpdatedAt, &characteristic.DeletedAt)

	if !reflect.DeepEqual(characteristic, want) {
		t.Errorf("got characteristic %+v, want %+v", characteristic, want)
	}
}

func TestCharacteristicsStore_Create_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	base_characteristic := newCharacteristic(t, tx)
	characteristic := models.Characteristic{base_characteristic, []int64(nil)}

	d := NewDatastore(tx)

	created, err := d.Characteristics.Create(&characteristic)
	if err != nil {
		t.Fatal(err)
	}
	if !created {
		t.Error("!created")
	}
	if characteristic.Id == 0 {
		t.Error("want nonzero characteristic.Id after submitting")
	}
}

func TestCharacteristicsStore_List_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	want_characteristic := insertCharacteristic(t, tx)
	want := []*models.Characteristic{want_characteristic}

	d := NewDatastore(tx)

	characteristics, err := d.Characteristics.List(&models.CharacteristicListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatal(err)
	}

	for i := range want {
		normalizeTime(&want[i].CreatedAt, &want[i].UpdatedAt, &want[i].DeletedAt)
		normalizeTime(&characteristics[i].CreatedAt, &characteristics[i].UpdatedAt, &characteristics[i].DeletedAt)
	}
	if !reflect.DeepEqual(characteristics, want) {
		t.Errorf("got characteristics %+v, want %+v", characteristics, want)
	}
}

func TestCharacteristicsStore_Update_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	characteristic := insertCharacteristic(t, tx)

	d := NewDatastore(tx)

	// Tweak it
	characteristic.CharacteristicName = "Updated Char"
	updated, err := d.Characteristics.Update(characteristic.Id, characteristic)
	if err != nil {
		t.Fatal(err)
	}

	if !updated {
		t.Error("!updated")
	}
}

func TestCharacteristicsStore_Delete_db(t *testing.T) {
	tx, _ := DB.Begin()
	defer tx.Rollback()

	characteristic := insertCharacteristic(t, tx)

	d := NewDatastore(tx)

	// Delete it
	deleted, err := d.Characteristics.Delete(characteristic.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !deleted {
		t.Error("!delete")
	}
}

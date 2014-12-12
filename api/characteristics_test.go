package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newCharacteristic() *models.Characteristic {
	characteristic := models.NewCharacteristic()
	return characteristic
}

func TestCharacteristic_Get(t *testing.T) {
	setup()

	want := newCharacteristic()

	calledGet := false

	store.Characteristics.(*models.MockCharacteristicsService).Get_ = func(id int64) (*models.Characteristic, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Characteristics.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got %+v but wanted %+v", got, want)
	}
}

func TestCharacteristic_Create(t *testing.T) {
	setup()

	want := newCharacteristic()

	calledPost := false
	store.Characteristics.(*models.MockCharacteristicsService).Create_ = func(characteristic *models.Characteristic) (bool, error) {
		if !normalizeDeepEqual(want, characteristic) {
			t.Errorf("wanted request for characteristic %d but got %d", want, characteristic)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Characteristics.Create(want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPost {
		t.Error("!calledPost")
	}
	if !success {
		t.Error("!success")
	}
}

func TestCharacteristic_List(t *testing.T) {
	setup()

	want := []*models.Characteristic{newCharacteristic()}
	wantOpt := &models.CharacteristicListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Characteristics.(*models.MockCharacteristicsService).List_ = func(opt *models.CharacteristicListOptions) ([]*models.Characteristic, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	characteristics, err := apiClient.Characteristics.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &characteristics) {
		t.Errorf("got characteristics %+v but wanted characteristics %+v", characteristics, want)
	}
}

func TestCharacteristic_Update(t *testing.T) {
	setup()

	want := newCharacteristic()

	calledPut := false
	store.Characteristics.(*models.MockCharacteristicsService).Update_ = func(id int64, characteristic *models.Characteristic) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, characteristic) {
			t.Errorf("wanted request for characteristic %d but got %d", want, characteristic)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.Characteristics.Update(want.Id, want)
	if err != nil {
		t.Fatal(err)
	}

	if !calledPut {
		t.Error("!calledPut")
	}
	if !success {
		t.Error("!success")
	}
}

func TestCharacteristic_Delete(t *testing.T) {
	setup()

	want := newCharacteristic()

	calledDelete := false
	store.Characteristics.(*models.MockCharacteristicsService).Delete_ = func(id int64) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic %d but got %d", want.Id, id)
		}
		calledDelete = true
		return true, nil
	}

	success, err := apiClient.Characteristics.Delete(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledDelete {
		t.Error("!calledDelete")
	}
	if !success {
		t.Error("!success")
	}
}

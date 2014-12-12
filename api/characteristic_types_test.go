package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newCharacteristicType() *models.CharacteristicType {
	characteristic_type := models.NewCharacteristicType()
	return characteristic_type
}

func TestCharacteristicType_Get(t *testing.T) {
	setup()

	want := newCharacteristicType()

	calledGet := false

	store.CharacteristicTypes.(*models.MockCharacteristicTypesService).Get_ = func(id int64) (*models.CharacteristicType, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic_type %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.CharacteristicTypes.Get(want.Id)
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

func TestCharacteristicType_Create(t *testing.T) {
	setup()

	want := newCharacteristicType()

	calledPost := false
	store.CharacteristicTypes.(*models.MockCharacteristicTypesService).Create_ = func(characteristic_type *models.CharacteristicType) (bool, error) {
		if !normalizeDeepEqual(want, characteristic_type) {
			t.Errorf("wanted request for characteristic_type %d but got %d", want, characteristic_type)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.CharacteristicTypes.Create(want)
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

func TestCharacteristicType_List(t *testing.T) {
	setup()

	want := []*models.CharacteristicType{newCharacteristicType()}
	wantOpt := &models.CharacteristicTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.CharacteristicTypes.(*models.MockCharacteristicTypesService).List_ = func(opt *models.CharacteristicTypeListOptions) ([]*models.CharacteristicType, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	characteristic_types, err := apiClient.CharacteristicTypes.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &characteristic_types) {
		t.Errorf("got characteristic_types %+v but wanted characteristic_types %+v", characteristic_types, want)
	}
}

func TestCharacteristicType_Update(t *testing.T) {
	setup()

	want := newCharacteristicType()

	calledPut := false
	store.CharacteristicTypes.(*models.MockCharacteristicTypesService).Update_ = func(id int64, characteristic_type *models.CharacteristicType) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic_type %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, characteristic_type) {
			t.Errorf("wanted request for characteristic_type %d but got %d", want, characteristic_type)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.CharacteristicTypes.Update(want.Id, want)
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

func TestCharacteristicType_Delete(t *testing.T) {
	setup()

	want := newCharacteristicType()

	calledDelete := false
	store.CharacteristicTypes.(*models.MockCharacteristicTypesService).Delete_ = func(id int64) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for characteristic_type %d but got %d", want.Id, id)
		}
		calledDelete = true
		return true, nil
	}

	success, err := apiClient.CharacteristicTypes.Delete(want.Id)
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

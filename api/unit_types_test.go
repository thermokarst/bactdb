package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newUnitType() *models.UnitType {
	unit_type := models.NewUnitType()
	return unit_type
}

func TestUnitType_Get(t *testing.T) {
	setup()

	want := newUnitType()

	calledGet := false

	store.UnitTypes.(*models.MockUnitTypesService).Get_ = func(id int64) (*models.UnitType, error) {
		if id != want.Id {
			t.Errorf("wanted request for unit_type %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.UnitTypes.Get(want.Id)
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

func TestUnitType_Create(t *testing.T) {
	setup()

	want := newUnitType()

	calledPost := false
	store.UnitTypes.(*models.MockUnitTypesService).Create_ = func(unit_type *models.UnitType) (bool, error) {
		if !normalizeDeepEqual(want, unit_type) {
			t.Errorf("wanted request for unit_type %d but got %d", want, unit_type)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.UnitTypes.Create(want)
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

func TestUnitType_List(t *testing.T) {
	setup()

	want := []*models.UnitType{newUnitType()}
	wantOpt := &models.UnitTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.UnitTypes.(*models.MockUnitTypesService).List_ = func(opt *models.UnitTypeListOptions) ([]*models.UnitType, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	unit_types, err := apiClient.UnitTypes.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &unit_types) {
		t.Errorf("got unit_types %+v but wanted unit_types %+v", unit_types, want)
	}
}

func TestUnitType_Update(t *testing.T) {
	setup()

	want := newUnitType()

	calledPut := false
	store.UnitTypes.(*models.MockUnitTypesService).Update_ = func(id int64, unit_type *models.UnitType) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for unit_type %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, unit_type) {
			t.Errorf("wanted request for unit_type %d but got %d", want, unit_type)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.UnitTypes.Update(want.Id, want)
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

func TestUnitType_Delete(t *testing.T) {
	setup()

	want := newUnitType()

	calledDelete := false
	store.UnitTypes.(*models.MockUnitTypesService).Delete_ = func(id int64) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for unit_type %d but got %d", want.Id, id)
		}
		calledDelete = true
		return true, nil
	}

	success, err := apiClient.UnitTypes.Delete(want.Id)
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

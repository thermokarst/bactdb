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

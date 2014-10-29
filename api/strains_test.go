package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newStrain() *models.Strain {
	strain := models.NewStrain()
	strain.Id = 1
	strain.SpeciesId = 1
	return strain
}

func TestStrain_Get(t *testing.T) {
	setup()

	want := newStrain()

	calledGet := false

	store.Strains.(*models.MockStrainsService).Get_ = func(id int64) (*models.Strain, error) {
		if id != want.Id {
			t.Errorf("wanted request for strain %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Strains.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got strain %+v but wanted strain %+v", got, want)
	}
}

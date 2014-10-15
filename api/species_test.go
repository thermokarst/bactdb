package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func TestSpecies_Get(t *testing.T) {
	setup()

	want := &models.Species{Id: 1, GenusId: 1, SpeciesName: "Test Species"}

	calledGet := false
	store.Species.(*models.MockSpeciesService).Get_ = func(id int64) (*models.Species, error) {
		if id != want.Id {
			t.Errorf("wanted request for species %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Species.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got species %+v but wanted species %+v", got, want)
	}
}

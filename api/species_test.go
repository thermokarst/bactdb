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

func TestSpecies_Create(t *testing.T) {
	setup()

	want := &models.Species{Id: 1, GenusId: 1, SpeciesName: "Test Species"}

	calledPost := false
	store.Species.(*models.MockSpeciesService).Create_ = func(species *models.Species) (bool, error) {
		if !normalizeDeepEqual(want, species) {
			t.Errorf("wanted request for species %d but got %d", want, species)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Species.Create(want)
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

func TestSpecies_List(t *testing.T) {
	setup()

	want := []*models.Species{{Id: 1, GenusId: 1, SpeciesName: "Test Species"}}
	wantOpt := &models.SpeciesListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Species.(*models.MockSpeciesService).List_ = func(opt *models.SpeciesListOptions) ([]*models.Species, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	species, err := apiClient.Species.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &species) {
		t.Errorf("got species %+v but wanted species %+v", species, want)
	}
}

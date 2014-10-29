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

func TestStrain_Create(t *testing.T) {
	setup()

	want := newStrain()

	calledPost := false
	store.Strains.(*models.MockStrainsService).Create_ = func(strain *models.Strain) (bool, error) {
		if !normalizeDeepEqual(want, strain) {
			t.Errorf("wanted request for strain %d but got %d", want, strain)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Strains.Create(want)
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

func TestStrain_List(t *testing.T) {
	setup()

	want := []*models.Strain{newStrain()}
	wantOpt := &models.StrainListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Strains.(*models.MockStrainsService).List_ = func(opt *models.StrainListOptions) ([]*models.Strain, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	strains, err := apiClient.Strains.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &strains) {
		t.Errorf("got strains %+v but wanted strains %+v", strains, want)
	}
}

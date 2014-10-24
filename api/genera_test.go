package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newGenus() *models.Genus {
	genus := models.NewGenus()
	genus.Id = 1
	return genus
}

func TestGenus_Get(t *testing.T) {
	setup()

	want := newGenus()

	calledGet := false
	store.Genera.(*models.MockGeneraService).Get_ = func(id int64) (*models.Genus, error) {
		if id != want.Id {
			t.Errorf("wanted request for genus %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Genera.Get(want.Id)
	if err != nil {
		t.Fatal(err)
	}

	if !calledGet {
		t.Error("!calledGet")
	}
	if !normalizeDeepEqual(want, got) {
		t.Errorf("got genus %+v but wanted genus %+v", got, want)
	}
}

func TestGenus_Create(t *testing.T) {
	setup()

	want := newGenus()

	calledPost := false
	store.Genera.(*models.MockGeneraService).Create_ = func(genus *models.Genus) (bool, error) {
		if !normalizeDeepEqual(want, genus) {
			t.Errorf("wanted request for genus %d but got %d", want, genus)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Genera.Create(want)
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

func TestGenus_List(t *testing.T) {
	setup()

	want := []*models.Genus{newGenus()}
	wantOpt := &models.GenusListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Genera.(*models.MockGeneraService).List_ = func(opt *models.GenusListOptions) ([]*models.Genus, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	genera, err := apiClient.Genera.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}
	if !normalizeDeepEqual(&want, &genera) {
		t.Errorf("got genera %+v but wanted genera %+v", genera, want)
	}
}

func TestGenus_Update(t *testing.T) {
	setup()

	want := newGenus()

	calledPut := false
	store.Genera.(*models.MockGeneraService).Update_ = func(id int64, genus *models.Genus) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for genus %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, genus) {
			t.Errorf("wanted request for genus %d but got %d", want, genus)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.Genera.Update(1, want)
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

func TestGenus_Delete(t *testing.T) {
	setup()

	want := newGenus()

	calledDelete := false
	store.Genera.(*models.MockGeneraService).Delete_ = func(id int64) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for genus %d but got %d", want.Id, id)
		}
		calledDelete = true
		return true, nil
	}

	success, err := apiClient.Genera.Delete(1)
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

package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newObservationType() *models.ObservationType {
	observation_type := models.NewObservationType()
	return observation_type
}

func TestObservationType_Get(t *testing.T) {
	setup()

	want := newObservationType()

	calledGet := false

	store.ObservationTypes.(*models.MockObservationTypesService).Get_ = func(id int64) (*models.ObservationType, error) {
		if id != want.Id {
			t.Errorf("wanted request for observation_type %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.ObservationTypes.Get(want.Id)
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

func TestObservationType_Create(t *testing.T) {
	setup()

	want := newObservationType()

	calledPost := false
	store.ObservationTypes.(*models.MockObservationTypesService).Create_ = func(observation_type *models.ObservationType) (bool, error) {
		if !normalizeDeepEqual(want, observation_type) {
			t.Errorf("wanted request for observation_type %d but got %d", want, observation_type)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.ObservationTypes.Create(want)
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

func TestObservationType_List(t *testing.T) {
	setup()

	want := []*models.ObservationType{newObservationType()}
	wantOpt := &models.ObservationTypeListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.ObservationTypes.(*models.MockObservationTypesService).List_ = func(opt *models.ObservationTypeListOptions) ([]*models.ObservationType, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	observation_types, err := apiClient.ObservationTypes.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &observation_types) {
		t.Errorf("got observation_types %+v but wanted observation_types %+v", observation_types, want)
	}
}

func TestObservationType_Update(t *testing.T) {
	setup()

	want := newObservationType()

	calledPut := false
	store.ObservationTypes.(*models.MockObservationTypesService).Update_ = func(id int64, observation_type *models.ObservationType) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for observation_type %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, observation_type) {
			t.Errorf("wanted request for observation_type %d but got %d", want, observation_type)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.ObservationTypes.Update(want.Id, want)
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

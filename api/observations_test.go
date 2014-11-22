package api

import (
	"testing"

	"github.com/thermokarst/bactdb/models"
)

func newObservation() *models.Observation {
	observation := models.NewObservation()
	return observation
}

func TestObservation_Get(t *testing.T) {
	setup()

	want := newObservation()

	calledGet := false

	store.Observations.(*models.MockObservationsService).Get_ = func(id int64) (*models.Observation, error) {
		if id != want.Id {
			t.Errorf("wanted request for observation %d but got %d", want.Id, id)
		}
		calledGet = true
		return want, nil
	}

	got, err := apiClient.Observations.Get(want.Id)
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

func TestObservation_Create(t *testing.T) {
	setup()

	want := newObservation()

	calledPost := false
	store.Observations.(*models.MockObservationsService).Create_ = func(observation *models.Observation) (bool, error) {
		if !normalizeDeepEqual(want, observation) {
			t.Errorf("wanted request for observation %d but got %d", want, observation)
		}
		calledPost = true
		return true, nil
	}

	success, err := apiClient.Observations.Create(want)
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

func TestObservation_List(t *testing.T) {
	setup()

	want := []*models.Observation{newObservation()}
	wantOpt := &models.ObservationListOptions{ListOptions: models.ListOptions{Page: 1, PerPage: 10}}

	calledList := false
	store.Observations.(*models.MockObservationsService).List_ = func(opt *models.ObservationListOptions) ([]*models.Observation, error) {
		if !normalizeDeepEqual(wantOpt, opt) {
			t.Errorf("wanted options %d but got %d", wantOpt, opt)
		}
		calledList = true
		return want, nil
	}

	observations, err := apiClient.Observations.List(wantOpt)
	if err != nil {
		t.Fatal(err)
	}

	if !calledList {
		t.Error("!calledList")
	}

	if !normalizeDeepEqual(&want, &observations) {
		t.Errorf("got observations %+v but wanted observations %+v", observations, want)
	}
}

func TestObservation_Update(t *testing.T) {
	setup()

	want := newObservation()

	calledPut := false
	store.Observations.(*models.MockObservationsService).Update_ = func(id int64, observation *models.Observation) (bool, error) {
		if id != want.Id {
			t.Errorf("wanted request for observation %d but got %d", want.Id, id)
		}
		if !normalizeDeepEqual(want, observation) {
			t.Errorf("wanted request for observation %d but got %d", want, observation)
		}
		calledPut = true
		return true, nil
	}

	success, err := apiClient.Observations.Update(want.Id, want)
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

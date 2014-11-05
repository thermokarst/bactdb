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

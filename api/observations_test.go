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

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveObservation(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	observation, err := store.Observations.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, observation)
}

func serveCreateObservation(w http.ResponseWriter, r *http.Request) error {
	var observation models.Observation
	err := json.NewDecoder(r.Body).Decode(&observation)
	if err != nil {
		return err
	}

	created, err := store.Observations.Create(&observation)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, observation)
}

func serveObservationList(w http.ResponseWriter, r *http.Request) error {
	var opt models.ObservationListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	observations, err := store.Observations.List(&opt)
	if err != nil {
		return err
	}
	if observations == nil {
		observations = []*models.Observation{}
	}

	return writeJSON(w, observations)
}

func serveUpdateObservation(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var observation models.Observation
	err := json.NewDecoder(r.Body).Decode(&observation)
	if err != nil {
		return err
	}

	updated, err := store.Observations.Update(id, &observation)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, observation)
}

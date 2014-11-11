package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveObservationType(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	observation_type, err := store.ObservationTypes.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, observation_type)
}

func serveCreateObservationType(w http.ResponseWriter, r *http.Request) error {
	var observation_type models.ObservationType
	err := json.NewDecoder(r.Body).Decode(&observation_type)
	if err != nil {
		return err
	}

	created, err := store.ObservationTypes.Create(&observation_type)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, observation_type)
}

func serveObservationTypeList(w http.ResponseWriter, r *http.Request) error {
	var opt models.ObservationTypeListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	observation_types, err := store.ObservationTypes.List(&opt)
	if err != nil {
		return err
	}
	if observation_types == nil {
		observation_types = []*models.ObservationType{}
	}

	return writeJSON(w, observation_types)
}

func serveUpdateObservationType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var observation_type models.ObservationType
	err := json.NewDecoder(r.Body).Decode(&observation_type)
	if err != nil {
		return err
	}

	updated, err := store.ObservationTypes.Update(id, &observation_type)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, observation_type)
}

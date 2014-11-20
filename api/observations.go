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

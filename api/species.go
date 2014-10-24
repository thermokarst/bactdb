package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveSpecies(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	species, err := store.Species.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, species)
}

func serveCreateSpecies(w http.ResponseWriter, r *http.Request) error {
	var species models.Species
	err := json.NewDecoder(r.Body).Decode(&species)
	if err != nil {
		return err
	}

	created, err := store.Species.Create(&species)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, species)
}

func serveSpeciesList(w http.ResponseWriter, r *http.Request) error {
	var opt models.SpeciesListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	species, err := store.Species.List(&opt)
	if err != nil {
		return err
	}
	if species == nil {
		species = []*models.Species{}
	}

	return writeJSON(w, species)
}

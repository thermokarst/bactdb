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

	return writeJSON(w, models.SpeciesJSON{species})
}

func serveCreateSpecies(w http.ResponseWriter, r *http.Request) error {
	var species models.SpeciesJSON
	err := json.NewDecoder(r.Body).Decode(&species)
	if err != nil {
		return err
	}

	created, err := store.Species.Create(species.Species)
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

	return writeJSON(w, models.SpeciesListJSON{Species: species})
}

func serveUpdateSpecies(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var species models.SpeciesJSON
	err := json.NewDecoder(r.Body).Decode(&species)
	if err != nil {
		return err
	}

	updated, err := store.Species.Update(id, species.Species)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, species)
}

func serveDeleteSpecies(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.Species.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.Species{})
}

func serveSubrouterSpeciesList(w http.ResponseWriter, r *http.Request) error {
	var opt models.SpeciesListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	opt.Genus = mux.Vars(r)["genus"]

	species, err := store.Species.List(&opt)
	if err != nil {
		return err
	}
	if species == nil {
		species = []*models.Species{}
	}

	return writeJSON(w, models.SpeciesListJSON{Species: species})
}

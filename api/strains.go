package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveStrain(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	strain, err := store.Strains.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, strain)
}

func serveCreateStrain(w http.ResponseWriter, r *http.Request) error {
	var strain models.Strain
	err := json.NewDecoder(r.Body).Decode(&strain)
	if err != nil {
		return err
	}

	created, err := store.Strains.Create(&strain)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, strain)
}

func serveStrainList(w http.ResponseWriter, r *http.Request) error {
	var opt models.StrainListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	strains, err := store.Strains.List(&opt)
	if err != nil {
		return err
	}
	if strains == nil {
		strains = []*models.Strain{}
	}

	return writeJSON(w, strains)
}

func serveUpdateStrain(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var strain models.Strain
	err := json.NewDecoder(r.Body).Decode(&strain)
	if err != nil {
		return err
	}

	updated, err := store.Strains.Update(id, &strain)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, strain)
}

func serveDeleteStrain(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.Strains.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.Strain{})
}

func serveSubrouterStrainsList(w http.ResponseWriter, r *http.Request) error {
	var opt models.StrainListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	opt.Genus = mux.Vars(r)["genus"]

	strains, err := store.Strains.List(&opt)
	if err != nil {
		return err
	}
	if strains == nil {
		strains = []*models.Strain{}
	}

	return writeJSON(w, strains)
}

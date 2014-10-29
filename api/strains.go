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

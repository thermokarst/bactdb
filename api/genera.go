package api

import (
	"encoding/json"
	"strconv"

	"github.com/gorilla/mux"

	"net/http"

	"github.com/thermokarst/bactdb/models"
)

func serveGenus(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	genus, err := store.Genera.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, models.GenusJSON{Genus: genus})
}

func serveCreateGenus(w http.ResponseWriter, r *http.Request) error {
	var genus models.GenusJSON
	err := json.NewDecoder(r.Body).Decode(&genus)
	if err != nil {
		return err
	}

	created, err := store.Genera.Create(genus.Genus)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, genus)
}

func serveGenera(w http.ResponseWriter, r *http.Request) error {
	var opt models.GenusListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	genera, err := store.Genera.List(&opt)
	if err != nil {
		return err
	}
	if genera == nil {
		genera = []*models.Genus{}
	}

	return writeJSON(w, models.GeneraJSON{Genera: genera})
}

func serveUpdateGenus(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var genus models.GenusJSON
	err := json.NewDecoder(r.Body).Decode(&genus)
	if err != nil {
		return err
	}

	updated, err := store.Genera.Update(id, genus.Genus)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, genus)
}

func serveDeleteGenus(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.Genera.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, nil)
}

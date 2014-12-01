package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveUnitType(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	unit_type, err := store.UnitTypes.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, unit_type)
}

func serveCreateUnitType(w http.ResponseWriter, r *http.Request) error {
	var unit_type models.UnitType
	err := json.NewDecoder(r.Body).Decode(&unit_type)
	if err != nil {
		return err
	}

	created, err := store.UnitTypes.Create(&unit_type)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, unit_type)
}

func serveUnitTypeList(w http.ResponseWriter, r *http.Request) error {
	var opt models.UnitTypeListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	unit_types, err := store.UnitTypes.List(&opt)
	if err != nil {
		return err
	}
	if unit_types == nil {
		unit_types = []*models.UnitType{}
	}

	return writeJSON(w, unit_types)
}

func serveUpdateUnitType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var unit_type models.UnitType
	err := json.NewDecoder(r.Body).Decode(&unit_type)
	if err != nil {
		return err
	}

	updated, err := store.UnitTypes.Update(id, &unit_type)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, unit_type)
}

func serveDeleteUnitType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.UnitTypes.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.UnitType{})
}

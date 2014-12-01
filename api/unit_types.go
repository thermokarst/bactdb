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

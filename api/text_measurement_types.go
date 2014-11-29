package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveTextMeasurementType(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	text_measurement_type, err := store.TextMeasurementTypes.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, text_measurement_type)
}

func serveCreateTextMeasurementType(w http.ResponseWriter, r *http.Request) error {
	var text_measurement_type models.TextMeasurementType
	err := json.NewDecoder(r.Body).Decode(&text_measurement_type)
	if err != nil {
		return err
	}

	created, err := store.TextMeasurementTypes.Create(&text_measurement_type)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, text_measurement_type)
}

func serveTextMeasurementTypeList(w http.ResponseWriter, r *http.Request) error {
	var opt models.TextMeasurementTypeListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	text_measurement_types, err := store.TextMeasurementTypes.List(&opt)
	if err != nil {
		return err
	}
	if text_measurement_types == nil {
		text_measurement_types = []*models.TextMeasurementType{}
	}

	return writeJSON(w, text_measurement_types)
}

func serveUpdateTextMeasurementType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var text_measurement_type models.TextMeasurementType
	err := json.NewDecoder(r.Body).Decode(&text_measurement_type)
	if err != nil {
		return err
	}

	updated, err := store.TextMeasurementTypes.Update(id, &text_measurement_type)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, text_measurement_type)
}

func serveDeleteTextMeasurementType(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.TextMeasurementTypes.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, &models.TextMeasurementType{})
}
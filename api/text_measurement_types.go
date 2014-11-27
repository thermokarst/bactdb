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

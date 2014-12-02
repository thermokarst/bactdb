package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/models"
)

func serveMeasurement(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	measurement, err := store.Measurements.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, measurement)
}

func serveCreateMeasurement(w http.ResponseWriter, r *http.Request) error {
	var measurement models.Measurement
	err := json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		return err
	}

	created, err := store.Measurements.Create(&measurement)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
	}

	return writeJSON(w, measurement)
}

func serveMeasurementList(w http.ResponseWriter, r *http.Request) error {
	var opt models.MeasurementListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	measurements, err := store.Measurements.List(&opt)
	if err != nil {
		return err
	}
	if measurements == nil {
		measurements = []*models.Measurement{}
	}

	return writeJSON(w, measurements)
}

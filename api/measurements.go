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

	return writeJSON(w, models.MeasurementJSON{Measurement: measurement})
}

func serveCreateMeasurement(w http.ResponseWriter, r *http.Request) error {
	var measurement models.MeasurementJSON
	err := json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		return err
	}

	created, err := store.Measurements.Create(measurement.Measurement)
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

	return writeJSON(w, models.MeasurementsJSON{Measurements: measurements})
}

func serveUpdateMeasurement(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	var measurement models.MeasurementJSON
	err := json.NewDecoder(r.Body).Decode(&measurement)
	if err != nil {
		return err
	}

	updated, err := store.Measurements.Update(id, measurement.Measurement)
	if err != nil {
		return err
	}
	if updated {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, measurement)
}

func serveDeleteMeasurement(w http.ResponseWriter, r *http.Request) error {
	id, _ := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)

	deleted, err := store.Measurements.Delete(id)
	if err != nil {
		return err
	}
	if deleted {
		w.WriteHeader(http.StatusOK)
	}

	return writeJSON(w, nil)
}

func serveSubrouterMeasurementsList(w http.ResponseWriter, r *http.Request) error {
	var opt models.MeasurementListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	opt.Genus = mux.Vars(r)["genus"]

	measurements, err := store.Measurements.List(&opt)
	if err != nil {
		return err
	}
	if measurements == nil {
		measurements = []*models.Measurement{}
	}

	return writeJSON(w, models.MeasurementsJSON{Measurements: measurements})
}

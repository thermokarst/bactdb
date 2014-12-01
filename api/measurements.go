package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

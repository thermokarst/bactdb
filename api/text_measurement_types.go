package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

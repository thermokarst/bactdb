package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveObservation(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	observation, err := store.Observations.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, observation)
}

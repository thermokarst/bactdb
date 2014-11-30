package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

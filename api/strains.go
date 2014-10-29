package api

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func serveStrain(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	strain, err := store.Strains.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, strain)
}

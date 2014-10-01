package api

import (
	"strconv"

	"github.com/gorilla/mux"

	"net/http"

	"github.com/thermokarst/bactdb/models"
)

func serveUser(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["ID"], 10, 0)
	if err != nil {
		return err
	}

	user, err := store.Users.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, user)
}

func serveUsers(w http.ResponseWriter, r *http.Request) error {
	var opt models.UserListOptions
	if err := schemaDecoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	users, err := store.Users.List(&opt)
	if err != nil {
		return err
	}
	if users == nil {
		users = []*models.User{}
	}

	return writeJSON(w, users)
}

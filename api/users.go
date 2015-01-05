package api

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"net/http"

	"github.com/thermokarst/bactdb/models"
)

func serveUser(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
	if err != nil {
		return err
	}

	user, err := store.Users.Get(id)
	if err != nil {
		return err
	}

	return writeJSON(w, user)
}

func serveCreateUser(w http.ResponseWriter, r *http.Request) error {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}

	created, err := store.Users.Create(&user)
	if err != nil {
		return err
	}
	if created {
		w.WriteHeader(http.StatusCreated)
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

func serveAuthenticateUser(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")

	auth_level, err := store.Users.Authenticate(username, password)
	if err != nil {
		return err
	}

	t := jwt.New(jwt.GetSigningMethod("RS256"))
	t.Claims["AccessToken"] = auth_level
	t.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := t.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:       tokenName,
		Value:      tokenString,
		Path:       "/",
		RawExpires: "0",
	})

	return writeJSON(w, auth_level)
}

package api

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const (
	tokenName = "AccessToken"
)

var (
	verifyKey, signKey    []byte
	errWhileSigningToken  = errors.New("error while signing token")
	errPleaseLogIn        = errors.New("please log in")
	errWhileParsingCookie = errors.New("error while parsing cookie")
	errTokenExpired       = errors.New("token expired")
	errGenericError       = errors.New("generic error")
	errAccessDenied       = errors.New("insufficient privileges")
)

func SetupCerts() error {
	signkey := os.Getenv("PRIVATE_KEY")
	if signkey == "" {
		return errors.New("please set PRIVATE_KEY")
	}
	signKey = []byte(signkey)

	verifykey := os.Getenv("PUBLIC_KEY")
	if verifykey == "" {
		return errors.New("please set PUBLIC_KEY")
	}
	verifyKey = []byte(verifykey)

	return nil
}

type authHandler func(http.ResponseWriter, *http.Request) error

// Only accessible with a valid token
func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Even though writeJSON sets the content type, we need to set it here because
	// calls to WriteHeader write out the entire header.
	w.Header().Set("content-type", "application/json; charset=utf-8")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, Error{errPleaseLogIn})
		return
	}
	s := strings.Split(authHeader, " ")

	// Validate the token
	token, err := jwt.Parse(s[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(verifyKey), nil
	})

	// Branch out into the possible error from signing
	switch err.(type) {
	case nil: // No error
		if !token.Valid { // But may still be invalid
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, Error{errPleaseLogIn})
			return
		}
	case *jwt.ValidationError: // Something was wrong during the validation
		vErr := err.(*jwt.ValidationError)
		switch vErr.Errors {
		case jwt.ValidationErrorExpired:
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, Error{errTokenExpired})
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			writeJSON(w, Error{errGenericError})
			return
		}
	default: // Something else went wrong
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{errGenericError})
		return
	}
	genus := mux.Vars(r)["genus"]
	// We don't care about this if we aren't accessing one of the subrouter routes.
	if genus != "" && genus != token.Claims["genus"] {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{errAccessDenied})
		return
	}
	hErr := h(w, r)
	if hErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{hErr})
	}
}

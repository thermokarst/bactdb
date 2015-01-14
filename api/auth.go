package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func SetupCerts(p string) error {
	var err error
	if err != nil {
		log.Fatalf("Path error: ", err)
	}

	// openssl genrsa -out app.rsa keysize
	privKeyPath := fmt.Sprintf("%vapp.rsa", p)
	signKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatalf("Error reading private key: ", err)
		return err
	}

	// openssl rsa -in app.rsa -pubout > app.rsa.pub
	pubKeyPath := fmt.Sprintf("%vapp.rsa.pub", p)
	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatalf("Error reading public key: ", err)
		return err
	}
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
		return verifyKey, nil
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

package api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	tokenName   = "AccessToken"
)

var (
	verifyKey, signKey    []byte
	errWhileSigningToken  = errors.New("error while signing token")
	errPleaseLogIn        = errors.New("please log in")
	errWhileParsingCookie = errors.New("error while parsing cookie")
	errTokenExpired       = errors.New("token expired")
	errGenericError       = errors.New("generic error")
)

func init() {
	var err error
	dir, _ := filepath.Abs(filepath.Dir(privKeyPath))
	fmt.Println(dir)

	signKey, err = ioutil.ReadFile(privKeyPath)

	if err != nil {
		// Before exploding, check up one level...
		signKey, err = ioutil.ReadFile("../../" + privKeyPath)
		if err != nil {
			log.Fatalf("Error reading private key: ", err)
			return
		}
	}

	verifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		// Before exploding, check up one level...
		verifyKey, err = ioutil.ReadFile("../../" + pubKeyPath)
		if err != nil {
			log.Fatalf("Error reading public key: ", err)
			return
		}
	}
}

type authHandler func(http.ResponseWriter, *http.Request) error

// Only accessible with a valid token
func (h authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Even though writeJSON sets the content type, we need to set it here because
	// calls to WriteHeader write out the entire header.
	w.Header().Set("content-type", "application/json; charset=utf-8")
	tokenCookie, err := r.Cookie(tokenName)
	switch {
	case err == http.ErrNoCookie:
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, Error{errPleaseLogIn})
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{errWhileParsingCookie})
		return
	}

	if tokenCookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, Error{errPleaseLogIn})
		return
	}

	// Validate the token
	token, err := jwt.Parse(tokenCookie.Value, func(token *jwt.Token) (interface{}, error) {
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
	hErr := h(w, r)
	if hErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{hErr})
	}
}

func restrictedHandler(w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, Message{"great success"})
}

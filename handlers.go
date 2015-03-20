package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func Handler() http.Handler {
	m := mux.NewRouter()

	// Non-auth routes
	m.HandleFunc("/authenticate", serveAuthenticateUser).Methods("POST")

	// Path-based pattern matching subrouter
	s := m.PathPrefix("/{genus}").Subrouter()

	// Strains
	s.Handle("/strains", authHandler(serveStrainsList)).Methods("GET")
	s.Handle("/strains/{Id:.+}", authHandler(serveStrain)).Methods("GET")

	// Measurements
	s.Handle("/measurements", authHandler(serveMeasurementsList)).Methods("GET")
	s.Handle("/measurements/{Id:.+}", authHandler(serveMeasurement)).Methods("GET")

	return corsHandler(m)
}

func corsHandler(h http.Handler) http.Handler {
	cors := func(w http.ResponseWriter, r *http.Request) {
		domains := os.Getenv("DOMAINS")
		allowedDomains := strings.Split(domains, ",")
		if origin := r.Header.Get("Origin"); origin != "" {
			for _, s := range allowedDomains {
				if s == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				}
			}
		}
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(cors)
}

// Only accessible with a valid token
func authHandler(f func(http.ResponseWriter, *http.Request)) http.Handler {
	h := http.HandlerFunc(f)
	auth := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, errPleaseLogIn.Error(), http.StatusUnauthorized)
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
				http.Error(w, errPleaseLogIn.Error(), http.StatusUnauthorized)
				return
			}
		case *jwt.ValidationError: // Something was wrong during the validation
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				http.Error(w, errTokenExpired.Error(), http.StatusUnauthorized)
				return
			default:
				http.Error(w, errGenericError.Error(), http.StatusInternalServerError)
				return
			}
		default: // Something else went wrong
			http.Error(w, errGenericError.Error(), http.StatusInternalServerError)
			return
		}
		genus := mux.Vars(r)["genus"]
		// We don't care about this if we aren't accessing one of the subrouter routes.
		if genus != "" && genus != token.Claims["genus"] {
			http.Error(w, errAccessDenied.Error(), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(auth)
}

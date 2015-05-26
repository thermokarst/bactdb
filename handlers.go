package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/thermokarst/jwt"
)

type Claims struct {
	Name string
	Iss  string
	Sub  int64
	Role string
	Iat  int64
	Exp  int64
}

func Handler() http.Handler {
	claimsFunc := func(email string) (map[string]interface{}, error) {
		currentTime := time.Now()
		user, err := dbGetUserByEmail(email)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"name": user.Name,
			"iss":  "bactdb",
			"sub":  user.Id,
			"role": user.Role,
			"iat":  currentTime.Unix(),
			"exp":  currentTime.Add(time.Minute * 60 * 24).Unix(),
		}, nil
	}

	verifyClaims := func(claims []byte, r *http.Request) error {
		currentTime := time.Now()
		var c Claims
		err := json.Unmarshal(claims, &c)
		if err != nil {
			return err
		}
		if currentTime.After(time.Unix(c.Exp, 0)) {
			return errors.New("this token has expired")
		}
		context.Set(r, "claims", c)
		return nil
	}

	config := &jwt.Config{
		Secret: os.Getenv("SECRET"),
		Auth:   dbAuthenticate,
		Claims: claimsFunc,
	}

	j, err := jwt.New(config)
	if err != nil {
		panic(err)
	}

	m := mux.NewRouter()

	// Non-auth routes
	m.Handle("/authenticate", tokenHandler(j.GenerateToken())).Methods("POST")

	// Auth routes
	m.Handle("/users", j.Secure(http.HandlerFunc(serveUsersList), verifyClaims)).Methods("GET")
	m.Handle("/users/{Id:.+}", j.Secure(http.HandlerFunc(serveUser), verifyClaims)).Methods("GET")

	// Path-based pattern matching subrouter
	s := m.PathPrefix("/{genus}").Subrouter()

	type r struct {
		f http.HandlerFunc
		m string
		p string
	}

	routes := []r{
		r{serveStrainsList, "GET", "/strains"},
		r{serveStrain, "GET", "/strains/{Id:.+}"},
		r{serveUpdateStrain, "PUT", "/strains/{Id:.+}"},
		r{serveMeasurementsList, "GET", "/measurements"},
		r{serveMeasurement, "GET", "/measurements/{Id:.+}"},
		r{serveCharacteristicsList, "GET", "/characteristics"},
		r{serveCharacteristic, "GET", "/characteristics/{Id:.+}"},
	}

	for _, route := range routes {
		s.Handle(route.p, j.Secure(http.HandlerFunc(route.f), verifyClaims)).Methods(route.m)
	}

	return corsHandler(m)
}

func tokenHandler(h http.Handler) http.Handler {
	token := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Hackish, but we want the token in a JSON object
		w.Write([]byte(`{"token":"`))
		h.ServeHTTP(w, r)
		w.Write([]byte(`"}`))
	}
	return http.HandlerFunc(token)
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

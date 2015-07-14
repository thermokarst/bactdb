package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	Ref  string
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
			"ref":  "",
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
	userService := UserService{}
	strainService := StrainService{}
	speciesService := SpeciesService{}
	characteristicService := CharacteristicService{}
	measurementService := MeasurementService{}

	m.Handle("/authenticate", tokenHandler(j.GenerateToken())).Methods("POST")

	// Everything past here is lumped under a genus
	s := m.PathPrefix("/{genus}").Subrouter()

	s.Handle("/users", errorHandler(handleCreater(userService))).Methods("POST")
	s.Handle("/users/verify/{Nonce}", errorHandler(handleUserVerify)).Methods("GET")

	type r struct {
		f errorHandler
		m string
		p string
	}

	// Everything past this point requires a valid token
	routes := []r{
		r{handleLister(userService), "GET", "/users"},
		r{handleGetter(userService), "GET", "/users/{Id:.+}"},
		r{handleUpdater(userService), "PUT", "/users/{Id:.+}"},
		r{handleLister(speciesService), "GET", "/species"},
		r{handleCreater(speciesService), "POST", "/species"},
		r{handleGetter(speciesService), "GET", "/species/{Id:.+}"},
		r{handleUpdater(speciesService), "PUT", "/species/{Id:.+}"},
		r{handleLister(strainService), "GET", "/strains"},
		r{handleCreater(strainService), "POST", "/strains"},
		r{handleGetter(strainService), "GET", "/strains/{Id:.+}"},
		r{handleUpdater(strainService), "PUT", "/strains/{Id:.+}"},
		r{handleLister(characteristicService), "GET", "/characteristics"},
		r{handleGetter(characteristicService), "GET", "/characteristics/{Id:.+}"},
		r{handleLister(measurementService), "GET", "/measurements"},
		r{handleGetter(measurementService), "GET", "/measurements/{Id:.+}"},
	}

	for _, route := range routes {
		s.Handle(route.p, j.Secure(errorHandler(route.f), verifyClaims)).Methods(route.m)
	}

	return jsonHandler(corsHandler(m))
}

func handleGetter(g getter) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := getClaims(r)

		e, appErr := g.get(id, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleLister(l lister) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		opt := r.URL.Query()
		opt.Add("Genus", mux.Vars(r)["genus"])

		claims := getClaims(r)

		es, appErr := l.list(&opt, &claims)
		if appErr != nil {
			return appErr
		}
		data, err := es.marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleUpdater(u updater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		e, err := u.unmarshal(bodyBytes)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := getClaims(r)

		appErr := u.update(id, &e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleCreater(c creater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		e, err := c.unmarshal(bodyBytes)
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}

		claims := getClaims(r)

		appErr := c.create(&e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.marshal()
		if err != nil {
			return newJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func tokenHandler(h http.Handler) http.Handler {
	token := func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		h.ServeHTTP(recorder, r)

		for key, val := range recorder.Header() {
			w.Header()[key] = val
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(recorder.Code)

		tokenData := string(recorder.Body.Bytes())

		var data []byte

		if recorder.Code != 200 {
			data, _ = json.Marshal(struct {
				Error string `json:"error"`
			}{
				Error: tokenData,
			})
		} else {
			data, _ = json.Marshal(struct {
				Token string `json:"token"`
			}{
				Token: tokenData,
			})
		}

		w.Write(data)
		return

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

func jsonHandler(h http.Handler) http.Handler {
	j := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(j)
}

type errorHandler func(http.ResponseWriter, *http.Request) *appError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		w.WriteHeader(err.Status)
		fmt.Fprintln(w, err.Error.Error())
	}
}

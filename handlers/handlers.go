package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/nytimes/gziphandler"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/thermokarst/jwt"
	"github.com/thermokarst/bactdb/api"
	"github.com/thermokarst/bactdb/auth"
	"github.com/thermokarst/bactdb/errors"
	"github.com/thermokarst/bactdb/helpers"
	"github.com/thermokarst/bactdb/models"
	"github.com/thermokarst/bactdb/types"
)

func verifyClaims(claims []byte, r *http.Request) error {
	// TODO: use helper
	currentTime := time.Now()
	var c types.Claims
	err := json.Unmarshal(claims, &c)
	if err != nil {
		return err
	}
	if currentTime.After(time.Unix(c.Exp, 0)) {
		return errors.ExpiredToken
	}
	context.Set(r, "claims", c)
	return nil
}

func Handler() http.Handler {
	m := mux.NewRouter()
	userService := api.UserService{}
	strainService := api.StrainService{}
	speciesService := api.SpeciesService{}
	characteristicService := api.CharacteristicService{}
	measurementService := api.MeasurementService{}

	m.Handle("/authenticate", tokenHandler(auth.Middleware.Authenticate())).Methods("POST")
	m.Handle("/refresh", auth.Middleware.Secure(errorHandler(tokenRefresh(auth.Middleware)), verifyClaims)).Methods("POST")

	// Everything past here is lumped under a genus
	s := m.PathPrefix("/{genus}").Subrouter()

	s.Handle("/users", errorHandler(handleCreater(userService))).Methods("POST")
	s.Handle("/users/verify/{Nonce}", errorHandler(api.HandleUserVerify)).Methods("GET")
	s.Handle("/users/lockout", errorHandler(api.HandleUserLockout)).Methods("POST")

	s.Handle("/compare", auth.Middleware.Secure(errorHandler(api.HandleCompare), verifyClaims)).Methods("GET")

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
		r{handleCreater(characteristicService), "POST", "/characteristics"},
		r{handleGetter(characteristicService), "GET", "/characteristics/{Id:.+}"},
		r{handleUpdater(characteristicService), "PUT", "/characteristics/{Id:.+}"},
		r{handleLister(measurementService), "GET", "/measurements"},
		r{handleCreater(measurementService), "POST", "/measurements"},
		r{handleGetter(measurementService), "GET", "/measurements/{Id:.+}"},
		r{handleUpdater(measurementService), "PUT", "/measurements/{Id:.+}"},
		r{handleDeleter(measurementService), "DELETE", "/measurements/{Id:.+}"},
	}

	for _, route := range routes {
		s.Handle(route.p, auth.Middleware.Secure(errorHandler(route.f), verifyClaims)).Methods(route.m)
	}

	return jsonHandler(gziphandler.GzipHandler(corsHandler(m)))
}

func handleGetter(g api.Getter) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		e, appErr := g.Get(id, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleLister(l api.Lister) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		opt := r.URL.Query()
		opt.Add("Genus", mux.Vars(r)["genus"])

		claims := helpers.GetClaims(r)

		es, appErr := l.List(&opt, &claims)
		if appErr != nil {
			return appErr
		}
		data, err := es.Marshal()
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleUpdater(u api.Updater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		e, err := u.Unmarshal(bodyBytes)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := u.Update(id, &e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleCreater(c api.Creater) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		e, err := c.Unmarshal(bodyBytes)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := c.Create(&e, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

		data, err := e.Marshal()
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		w.Write(data)
		return nil
	}
}

func handleDeleter(d api.Deleter) errorHandler {
	return func(w http.ResponseWriter, r *http.Request) *types.AppError {
		id, err := strconv.ParseInt(mux.Vars(r)["Id"], 10, 0)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}

		claims := helpers.GetClaims(r)

		appErr := d.Delete(id, mux.Vars(r)["genus"], &claims)
		if appErr != nil {
			return appErr
		}

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

type errorHandler func(http.ResponseWriter, *http.Request) *types.AppError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		w.WriteHeader(err.Status)
		fmt.Fprintln(w, err.Error.Error())
	}
}

func tokenRefresh(j *jwt.Middleware) errorHandler {
	t := func(w http.ResponseWriter, r *http.Request) *types.AppError {
		claims := helpers.GetClaims(r)
		user, err := models.DbGetUserById(claims.Sub)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		user.Password = ""
		token, err := auth.Middleware.CreateToken(user.Email)
		if err != nil {
			return NewJSONError(err, http.StatusInternalServerError)
		}
		data, _ := json.Marshal(struct {
			Token string `json:"token"`
		}{
			Token: token,
		})
		w.Write(data)
		return nil
	}
	return t
}

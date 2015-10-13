package handlers

import (
	"net/http"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/nytimes/gziphandler"
	"github.com/thermokarst/bactdb/api"
	"github.com/thermokarst/bactdb/auth"
)

// Handler is the root HTTP handler for bactdb.
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
		r{api.HandleUserPasswordChange, "POST", "/users/password"},
		r{handleGetter(userService), "GET", "/users/{ID:.+}"},
		r{handleUpdater(userService), "PUT", "/users/{ID:.+}"},
		r{handleLister(speciesService), "GET", "/species"},
		r{handleCreater(speciesService), "POST", "/species"},
		r{handleGetter(speciesService), "GET", "/species/{ID:.+}"},
		r{handleUpdater(speciesService), "PUT", "/species/{ID:.+}"},
		r{handleDeleter(speciesService), "DELETE", "/species/{ID:.+}"},
		r{handleLister(strainService), "GET", "/strains"},
		r{handleCreater(strainService), "POST", "/strains"},
		r{handleGetter(strainService), "GET", "/strains/{ID:.+}"},
		r{handleUpdater(strainService), "PUT", "/strains/{ID:.+}"},
		r{handleDeleter(strainService), "DELETE", "/strains/{ID:.+}"},
		r{handleLister(characteristicService), "GET", "/characteristics"},
		r{handleCreater(characteristicService), "POST", "/characteristics"},
		r{handleGetter(characteristicService), "GET", "/characteristics/{ID:.+}"},
		r{handleUpdater(characteristicService), "PUT", "/characteristics/{ID:.+}"},
		r{handleDeleter(characteristicService), "DELETE", "/characteristics/{ID:.+}"},
		r{handleLister(measurementService), "GET", "/measurements"},
		r{handleCreater(measurementService), "POST", "/measurements"},
		r{handleGetter(measurementService), "GET", "/measurements/{ID:.+}"},
		r{handleUpdater(measurementService), "PUT", "/measurements/{ID:.+}"},
		r{handleDeleter(measurementService), "DELETE", "/measurements/{ID:.+}"},
	}

	for _, route := range routes {
		s.Handle(route.p, auth.Middleware.Secure(errorHandler(route.f), verifyClaims)).Methods(route.m)
	}

	return jsonHandler(gziphandler.GzipHandler(corsHandler(m)))
}

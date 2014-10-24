package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/thermokarst/bactdb/datastore"
	"github.com/thermokarst/bactdb/router"
)

var (
	store         = datastore.NewDatastore(nil)
	schemaDecoder = schema.NewDecoder()
)

func Handler() *mux.Router {
	m := router.API()

	m.Get(router.User).Handler(handler(serveUser))
	m.Get(router.CreateUser).Handler(handler(serveCreateUser))
	m.Get(router.Users).Handler(handler(serveUsers))

	m.Get(router.Genus).Handler(handler(serveGenus))
	m.Get(router.CreateGenus).Handler(handler(serveCreateGenus))
	m.Get(router.Genera).Handler(handler(serveGenera))
	m.Get(router.UpdateGenus).Handler(handler(serveUpdateGenus))
	m.Get(router.DeleteGenus).Handler(handler(serveDeleteGenus))

	m.Get(router.Species).Handler(handler(serveSpecies))
	m.Get(router.CreateSpecies).Handler(handler(serveCreateSpecies))
	m.Get(router.SpeciesList).Handler(handler(serveSpeciesList))
	m.Get(router.UpdateSpecies).Handler(handler(serveUpdateSpecies))

	return m
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s", err)
		log.Println(err)
	}
}

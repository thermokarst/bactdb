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
	m.Get(router.DeleteSpecies).Handler(handler(serveDeleteSpecies))

	m.Get(router.Strain).Handler(handler(serveStrain))
	m.Get(router.CreateStrain).Handler(handler(serveCreateStrain))
	m.Get(router.Strains).Handler(handler(serveStrainList))
	m.Get(router.UpdateStrain).Handler(handler(serveUpdateStrain))
	m.Get(router.DeleteStrain).Handler(handler(serveDeleteStrain))

	m.Get(router.ObservationType).Handler(handler(serveObservationType))
	m.Get(router.CreateObservationType).Handler(handler(serveCreateObservationType))
	m.Get(router.ObservationTypes).Handler(handler(serveObservationTypeList))
	m.Get(router.UpdateObservationType).Handler(handler(serveUpdateObservationType))
	m.Get(router.DeleteObservationType).Handler(handler(serveDeleteObservationType))

	m.Get(router.Observation).Handler(handler(serveObservation))
	m.Get(router.CreateObservation).Handler(handler(serveCreateObservation))
	m.Get(router.Observations).Handler(handler(serveObservationList))
	m.Get(router.UpdateObservation).Handler(handler(serveUpdateObservation))
	m.Get(router.DeleteObservation).Handler(handler(serveDeleteObservation))

	m.Get(router.TextMeasurementType).Handler(handler(serveTextMeasurementType))
	m.Get(router.CreateTextMeasurementType).Handler(handler(serveCreateTextMeasurementType))
	m.Get(router.TextMeasurementTypes).Handler(handler(serveTextMeasurementTypeList))
	m.Get(router.UpdateTextMeasurementType).Handler(handler(serveUpdateTextMeasurementType))
	m.Get(router.DeleteTextMeasurementType).Handler(handler(serveDeleteTextMeasurementType))

	m.Get(router.UnitType).Handler(handler(serveUnitType))
	m.Get(router.CreateUnitType).Handler(handler(serveCreateUnitType))
	m.Get(router.UnitTypes).Handler(handler(serveUnitTypeList))
	m.Get(router.UpdateUnitType).Handler(handler(serveUpdateUnitType))
	m.Get(router.DeleteUnitType).Handler(handler(serveDeleteUnitType))

	m.Get(router.Measurement).Handler(handler(serveMeasurement))
	m.Get(router.CreateMeasurement).Handler(handler(serveCreateMeasurement))
	m.Get(router.Measurements).Handler(handler(serveMeasurementList))
	m.Get(router.UpdateMeasurement).Handler(handler(serveUpdateMeasurement))

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

package api

import (
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

	m.Get(router.User).Handler(authHandler(serveUser))
	m.Get(router.CreateUser).Handler(authHandler(serveCreateUser))
	m.Get(router.Users).Handler(authHandler(serveUsers))
	m.Get(router.GetToken).Handler(handler(serveAuthenticateUser))

	m.Get(router.Genus).Handler(authHandler(serveGenus))
	m.Get(router.CreateGenus).Handler(authHandler(serveCreateGenus))
	m.Get(router.Genera).Handler(authHandler(serveGenera))
	m.Get(router.UpdateGenus).Handler(authHandler(serveUpdateGenus))
	m.Get(router.DeleteGenus).Handler(authHandler(serveDeleteGenus))

	m.Get(router.Species).Handler(authHandler(serveSpecies))
	m.Get(router.CreateSpecies).Handler(authHandler(serveCreateSpecies))
	m.Get(router.SpeciesList).Handler(authHandler(serveSpeciesList))
	m.Get(router.UpdateSpecies).Handler(authHandler(serveUpdateSpecies))
	m.Get(router.DeleteSpecies).Handler(authHandler(serveDeleteSpecies))

	m.Get(router.Strain).Handler(authHandler(serveStrain))
	m.Get(router.CreateStrain).Handler(authHandler(serveCreateStrain))
	m.Get(router.Strains).Handler(authHandler(serveStrainList))
	m.Get(router.UpdateStrain).Handler(authHandler(serveUpdateStrain))
	m.Get(router.DeleteStrain).Handler(authHandler(serveDeleteStrain))

	m.Get(router.CharacteristicType).Handler(authHandler(serveCharacteristicType))
	m.Get(router.CreateCharacteristicType).Handler(authHandler(serveCreateCharacteristicType))
	m.Get(router.CharacteristicTypes).Handler(authHandler(serveCharacteristicTypeList))
	m.Get(router.UpdateCharacteristicType).Handler(authHandler(serveUpdateCharacteristicType))
	m.Get(router.DeleteCharacteristicType).Handler(authHandler(serveDeleteCharacteristicType))

	m.Get(router.Characteristic).Handler(authHandler(serveCharacteristic))
	m.Get(router.CreateCharacteristic).Handler(authHandler(serveCreateCharacteristic))
	m.Get(router.Characteristics).Handler(authHandler(serveCharacteristicList))
	m.Get(router.UpdateCharacteristic).Handler(authHandler(serveUpdateCharacteristic))
	m.Get(router.DeleteCharacteristic).Handler(authHandler(serveDeleteCharacteristic))

	m.Get(router.TextMeasurementType).Handler(authHandler(serveTextMeasurementType))
	m.Get(router.CreateTextMeasurementType).Handler(authHandler(serveCreateTextMeasurementType))
	m.Get(router.TextMeasurementTypes).Handler(authHandler(serveTextMeasurementTypeList))
	m.Get(router.UpdateTextMeasurementType).Handler(authHandler(serveUpdateTextMeasurementType))
	m.Get(router.DeleteTextMeasurementType).Handler(authHandler(serveDeleteTextMeasurementType))

	m.Get(router.UnitType).Handler(authHandler(serveUnitType))
	m.Get(router.CreateUnitType).Handler(authHandler(serveCreateUnitType))
	m.Get(router.UnitTypes).Handler(authHandler(serveUnitTypeList))
	m.Get(router.UpdateUnitType).Handler(authHandler(serveUpdateUnitType))
	m.Get(router.DeleteUnitType).Handler(authHandler(serveDeleteUnitType))

	m.Get(router.Measurement).Handler(authHandler(serveMeasurement))
	m.Get(router.CreateMeasurement).Handler(authHandler(serveCreateMeasurement))
	m.Get(router.Measurements).Handler(authHandler(serveMeasurementList))
	m.Get(router.UpdateMeasurement).Handler(authHandler(serveUpdateMeasurement))
	m.Get(router.DeleteMeasurement).Handler(authHandler(serveDeleteMeasurement))

	m.Get(router.SubrouterListSpecies).Handler(authHandler(serveSubrouterSpeciesList))
	m.Get(router.SubrouterListStrains).Handler(authHandler(serveSubrouterStrainsList))
	m.Get(router.SubrouterListMeasurements).Handler(authHandler(serveSubrouterMeasurementsList))

	return m
}

type handler func(http.ResponseWriter, *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, Error{err})
	}
}

package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()

	// Users
	m.Path("/users").Methods("GET").Name(Users)
	m.Path("/users").Methods("POST").Name(CreateUser)
	m.Path("/users/{Id:.+}").Methods("GET").Name(User)
	m.Path("/authenticate/").Methods("POST").Name(GetToken)

	// Genera
	m.Path("/genera").Methods("GET").Name(Genera)
	m.Path("/genera").Methods("POST").Name(CreateGenus)
	m.Path("/genera/{Id:.+}").Methods("GET").Name(Genus)
	m.Path("/genera/{Id:.+}").Methods("PUT").Name(UpdateGenus)
	m.Path("/genera/{Id:.+}").Methods("DELETE").Name(DeleteGenus)

	// Species
	m.Path("/species").Methods("GET").Name(SpeciesList)
	m.Path("/species").Methods("POST").Name(CreateSpecies)
	m.Path("/species/{Id:.+}").Methods("GET").Name(Species)
	m.Path("/species/{Id:.+}").Methods("PUT").Name(UpdateSpecies)
	m.Path("/species/{Id:.+}").Methods("DELETE").Name(DeleteSpecies)

	// Strains
	m.Path("/strains").Methods("GET").Name(Strains)
	m.Path("/strains").Methods("POST").Name(CreateStrain)
	m.Path("/strains/{Id:.+}").Methods("GET").Name(Strain)
	m.Path("/strains/{Id:.+}").Methods("PUT").Name(UpdateStrain)
	m.Path("/strains/{Id:.+}").Methods("DELETE").Name(DeleteStrain)

	// CharacteristicTypes
	m.Path("/characteristic_types").Methods("GET").Name(CharacteristicTypes)
	m.Path("/characteristic_types").Methods("POST").Name(CreateCharacteristicType)
	m.Path("/characteristic_types/{Id:.+}").Methods("GET").Name(CharacteristicType)
	m.Path("/characteristic_types/{Id:.+}").Methods("PUT").Name(UpdateCharacteristicType)
	m.Path("/characteristic_types/{Id:.+}").Methods("DELETE").Name(DeleteCharacteristicType)

	// Characteristics
	m.Path("/characteristics").Methods("GET").Name(Characteristics)
	m.Path("/characteristics").Methods("POST").Name(CreateCharacteristic)
	m.Path("/characteristics/{Id:.+}").Methods("GET").Name(Characteristic)
	m.Path("/characteristics/{Id:.+}").Methods("PUT").Name(UpdateCharacteristic)
	m.Path("/characteristics/{Id:.+}").Methods("DELETE").Name(DeleteCharacteristic)

	// TextMeasurementTypes
	m.Path("/text_measurement_types/").Methods("GET").Name(TextMeasurementTypes)
	m.Path("/text_measurement_types/").Methods("POST").Name(CreateTextMeasurementType)
	m.Path("/text_measurement_types/{Id:.+}").Methods("GET").Name(TextMeasurementType)
	m.Path("/text_measurement_types/{Id:.+}").Methods("PUT").Name(UpdateTextMeasurementType)
	m.Path("/text_measurement_types/{Id:.+}").Methods("DELETE").Name(DeleteTextMeasurementType)

	// UnitTypes
	m.Path("/unit_types/").Methods("GET").Name(UnitTypes)
	m.Path("/unit_types/").Methods("POST").Name(CreateUnitType)
	m.Path("/unit_types/{Id:.+}").Methods("GET").Name(UnitType)
	m.Path("/unit_types/{Id:.+}").Methods("PUT").Name(UpdateUnitType)
	m.Path("/unit_types/{Id:.+}").Methods("DELETE").Name(DeleteUnitType)

	// Measurements
	m.Path("/measurements/").Methods("GET").Name(Measurements)
	m.Path("/measurements/").Methods("POST").Name(CreateMeasurement)
	m.Path("/measurements/{Id:.+}").Methods("GET").Name(Measurement)
	m.Path("/measurements/{Id:.+}").Methods("PUT").Name(UpdateMeasurement)
	m.Path("/measurements/{Id:.+}").Methods("DELETE").Name(DeleteMeasurement)

	// Subrouter for auth/security
	s := m.PathPrefix("/{genus}").Subrouter()
	s.Path("/species").Methods("GET").Name(SubrouterListSpecies)
	s.Path("/strains").Methods("GET").Name(SubrouterListStrains)
	s.Path("/measurements").Methods("GET").Name(SubrouterListMeasurements)

	return m
}

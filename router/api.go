package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	m := mux.NewRouter()

	// Users
	m.Path("/users").Methods("GET").Name(Users)
	m.Path("/users").Methods("POST").Name(CreateUser)
	m.Path("/users/{Id:.+}").Methods("GET").Name(User)

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

	// ObservationTypes
	m.Path("/observation_types").Methods("GET").Name(ObservationTypes)
	m.Path("/observation_types").Methods("POST").Name(CreateObservationType)
	m.Path("/observation_types/{Id:.+}").Methods("GET").Name(ObservationType)
	m.Path("/observation_types/{Id:.+}").Methods("PUT").Name(UpdateObservationType)
	m.Path("/observation_types/{Id:.+}").Methods("DELETE").Name(DeleteObservationType)

	// Observations
	m.Path("/observations").Methods("GET").Name(Observations)
	m.Path("/observations").Methods("POST").Name(CreateObservation)
	m.Path("/observations/{Id:.+}").Methods("GET").Name(Observation)
	m.Path("/observations/{Id:.+}").Methods("PUT").Name(UpdateObservation)
	m.Path("/observations/{Id:.+}").Methods("DELETE").Name(DeleteObservation)

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

	// Authentication
	m.Path("/token/").Methods("GET").Name(GetToken)
	m.Path("/restricted/").Methods("GET").Name(Restricted)

	return m
}

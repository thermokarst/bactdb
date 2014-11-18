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
	m.Path("/observations/{Id:.+}").Methods("GET").Name(Observation)

	return m
}

package router

const (
	User       = "users:get"
	CreateUser = "users:create"
	Users      = "users:list"

	Genus       = "genus:get"
	CreateGenus = "genus:create"
	Genera      = "genus:list"
	UpdateGenus = "genus:update"
	DeleteGenus = "genus:delete"

	Species       = "species:get"
	CreateSpecies = "species:create"
	SpeciesList   = "species:list"
	UpdateSpecies = "species:update"
	DeleteSpecies = "species:delete"

	Strain       = "strain:get"
	CreateStrain = "strain:create"
	Strains      = "strain:list"
	UpdateStrain = "strain:update"
	DeleteStrain = "strain:delete"

	ObservationType       = "observation_type:get"
	CreateObservationType = "observation_type:create"
	ObservationTypes      = "observation_type:list"
	UpdateObservationType = "observation_type:update"
	DeleteObservationType = "observation_type:delete"

	Observation       = "observation:get"
	CreateObservation = "observation:create"
)

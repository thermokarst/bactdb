package datastore

import (
	"errors"

	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

// A datastore access point (in PostgreSQL)
type Datastore struct {
	Users            models.UsersService
	Genera           models.GeneraService
	Species          models.SpeciesService
	Strains          models.StrainsService
	ObservationTypes models.ObservationTypesService
	dbh              modl.SqlExecutor
}

var (
	ErrNoRowsUpdated = errors.New(`no rows updated`)
	ErrNoRowsDeleted = errors.New(`no rows deleted`)
)

// NewDatastore creates a new client for accessing the datastore (in PostgreSQL).
// If dbh is nil, it uses the global DB handle.
func NewDatastore(dbh modl.SqlExecutor) *Datastore {
	if dbh == nil {
		dbh = DBH
	}

	d := &Datastore{dbh: dbh}
	d.Users = &usersStore{d}
	d.Genera = &generaStore{d}
	d.Species = &speciesStore{d}
	d.Strains = &strainsStore{d}
	d.ObservationTypes = &observationTypesStore{d}
	return d
}

func NewMockDatastore() *Datastore {
	return &Datastore{
		Users:            &models.MockUsersService{},
		Genera:           &models.MockGeneraService{},
		Species:          &models.MockSpeciesService{},
		Strains:          &models.MockStrainsService{},
		ObservationTypes: &models.MockObservationTypesService{},
	}
}

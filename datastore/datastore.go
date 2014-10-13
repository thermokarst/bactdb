package datastore

import (
	"github.com/jmoiron/modl"
	"github.com/thermokarst/bactdb/models"
)

// A datastore access point (in PostgreSQL)
type Datastore struct {
	Users  models.UsersService
	Genera models.GeneraService
	dbh    modl.SqlExecutor
}

// NewDatastore creates a new client for accessing the datastore (in PostgreSQL).
// If dbh is nil, it uses the global DB handle.
func NewDatastore(dbh modl.SqlExecutor) *Datastore {
	if dbh == nil {
		dbh = DBH
	}

	d := &Datastore{dbh: dbh}
	d.Users = &usersStore{d}
	d.Genera = &generaStore{d}
	return d
}

func NewMockDatastore() *Datastore {
	return &Datastore{
		Users:  &models.MockUsersService{},
		Genera: &models.MockGeneraService{},
	}
}

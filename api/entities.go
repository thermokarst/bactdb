package api

import (
	"net/url"

	"github.com/thermokarst/bactdb/types"
)

// Getter gets a single entity.
type Getter interface {
	Get(int64, string, *types.Claims) (types.Entity, *types.AppError)
}

// Lister lists entities.
type Lister interface {
	List(*url.Values, *types.Claims) (types.Entity, *types.AppError)
}

// Updater updates entities.
type Updater interface {
	Update(int64, *types.Entity, string, *types.Claims) *types.AppError
	Unmarshal([]byte) (types.Entity, error)
}

// Creater creates entities.
type Creater interface {
	Create(*types.Entity, string, *types.Claims) *types.AppError
	Unmarshal([]byte) (types.Entity, error)
}

// Deleter deletes entities.
type Deleter interface {
	Delete(int64, string, *types.Claims) *types.AppError
}

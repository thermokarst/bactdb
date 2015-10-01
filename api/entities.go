package api

import (
	"net/url"

	"github.com/thermokarst/bactdb/types"
)

type Getter interface {
	Get(int64, string, *types.Claims) (types.Entity, *types.AppError)
}

type Lister interface {
	List(*url.Values, *types.Claims) (types.Entity, *types.AppError)
}

type Updater interface {
	Update(int64, *types.Entity, string, *types.Claims) *types.AppError
	Unmarshal([]byte) (types.Entity, error)
}

type Creater interface {
	Create(*types.Entity, string, *types.Claims) *types.AppError
	Unmarshal([]byte) (types.Entity, error)
}
type Deleter interface {
	Delete(int64, string, *types.Claims) *types.AppError
}

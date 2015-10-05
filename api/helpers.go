package api

import "github.com/thermokarst/bactdb/types"

func newJSONError(err error, status int) *types.AppError {
	return &types.AppError{
		Error:  types.ErrorJSON{Err: err},
		Status: status,
	}
}

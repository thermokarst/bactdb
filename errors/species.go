package errors

import "errors"

var (
	SpeciesNotFound   = errors.New("Species not found")
	SpeciesNotUpdated = errors.New("Species not updated")
)

package errors

import "errors"

var (
	// ErrSpeciesNotFound when not found.
	ErrSpeciesNotFound = errors.New("Species not found")
	// ErrSpeciesNotUpdated when not updated.
	ErrSpeciesNotUpdated = errors.New("Species not updated")
	// ErrSpeciesNotDeleted when not deleted.
	ErrSpeciesNotDeleted = errors.New("Species not deleted")
)

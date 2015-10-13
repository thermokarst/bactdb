package errors

import "errors"

var (
	// ErrStrainNotFound when not found.
	ErrStrainNotFound = errors.New("Strain not found")
	// ErrStrainNotUpdated when not updated.
	ErrStrainNotUpdated = errors.New("Strain not updated")
	// ErrStrainNotDeleted when not deleted.
	ErrStrainNotDeleted = errors.New("Strain not deleted")
)

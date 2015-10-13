package errors

import "errors"

var (
	// ErrMeasurementNotFound when not found.
	ErrMeasurementNotFound = errors.New("Measurement not found")
	// ErrMeasurementNotUpdated when not updated.
	ErrMeasurementNotUpdated = errors.New("Measurement not updated")
	// ErrMeasurementNotDeleted when not deleted.
	ErrMeasurementNotDeleted = errors.New("Measurement not deleted")
)

package errors

import "errors"

var (
	// ErrMeasurementNotFound when not found.
	ErrMeasurementNotFound = errors.New("Measurement not found")
	// ErrMeasurementNotUpdate when not updated.
	ErrMeasurementNotUpdated = errors.New("Measurement not updated")
)

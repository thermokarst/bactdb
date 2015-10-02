package errors

import "errors"

var (
	// ErrMeasurementNotFound when not found.
	ErrMeasurementNotFound = errors.New("Measurement not found")
)

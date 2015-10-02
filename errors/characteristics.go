package errors

import "errors"

var (
	// ErrCharacteristicNotFound when not found.
	ErrCharacteristicNotFound = errors.New("Characteristic not found")
	// ErrCharacteristicNotUpdated when not updated.
	ErrCharacteristicNotUpdated = errors.New("Characteristic not updated")
)

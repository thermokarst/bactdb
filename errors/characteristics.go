package errors

import "errors"

var (
	CharacteristicNotFound   = errors.New("Characteristic not found")
	CharacteristicNotUpdated = errors.New("Characteristic not updated")
)

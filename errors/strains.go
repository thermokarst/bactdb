package errors

import "errors"

var (
	StrainNotFound   = errors.New("Strain not found")
	StrainNotUpdated = errors.New("Strain not updated")
)

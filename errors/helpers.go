package errors

import "errors"

var (
	// ErrMustProvideOptions when missing options.
	ErrMustProvideOptions = errors.New("Must provide necessary options")
)

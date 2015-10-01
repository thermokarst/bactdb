package errors

import "errors"

var (
	MustProvideOptions = errors.New("Must provide necessary options")
)

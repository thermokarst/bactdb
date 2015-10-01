package errors

import "errors"

var (
	ExpiredToken = errors.New("this token has expired")
)

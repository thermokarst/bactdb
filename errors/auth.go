package errors

import "errors"

var (
	// ErrExpiredToken when expired token.
	ErrExpiredToken = errors.New("this token has expired")
)

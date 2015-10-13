package errors

import "errors"

var (
	// ErrExpiredToken when expired token.
	ErrExpiredToken = errors.New("this token has expired")
	// ErrInvalidToken when the role doesn't match the DB
	ErrInvalidToken = errors.New("this token needs to be reissued")
)

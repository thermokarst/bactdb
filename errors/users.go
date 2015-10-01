package errors

import "errors"

var (
	UserNotFound           = errors.New("No user found")
	UserNotUpdated         = errors.New("Count 0")
	UserMissingEmail       = errors.New("Missing email")
	InvalidEmailOrPassword = errors.New("Invalid email or password")
	EmailAddressTaken      = errors.New("Email address is already registered")
)

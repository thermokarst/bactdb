package errors

import "errors"

var (
	// ErrUserNotFound when not found.
	ErrUserNotFound = errors.New("No user found")
	// ErrUserNotUpdated when not updated.
	ErrUserNotUpdated = errors.New("User not updated")
	// ErrUserNotDeleted when not deleted.
	ErrUserNotDeleted = errors.New("User not deleted")
	// ErrUserMissingEmail when missing email.
	ErrUserMissingEmail = errors.New("Missing email")
	// ErrInvalidEmailOrPassword when invalid login credentials.
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
	// ErrEmailAddressTaken when email already registered.
	ErrEmailAddressTaken = errors.New("Email address is already registered")
	// ErrUserForbidden when user not allowed to view a resource
	ErrUserForbidden = errors.New("User ccount not authorized")
)

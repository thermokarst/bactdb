package types

import "encoding/json"

// ErrorJSON is an error that serializes to a JSON-encoded representation of the
// error message.
type ErrorJSON struct {
	Err error
}

// Error satisfies the necessary interface to make ErrorJSON an error.
func (ej ErrorJSON) Error() string {
	e, _ := json.Marshal(struct {
		Err string `json:"error"`
	}{
		Err: ej.Err.Error(),
	})
	return string(e)
}

// AppError returns an error plus an HTTP status code.
type AppError struct {
	Error  error
	Status int
}

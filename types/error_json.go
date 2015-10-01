package types

import "encoding/json"

type ErrorJSON struct {
	Err error
}

func (ej ErrorJSON) Error() string {
	e, _ := json.Marshal(struct {
		Err string `json:"error"`
	}{
		Err: ej.Err.Error(),
	})
	return string(e)
}

type AppError struct {
	Error  error
	Status int
}

func NewJSONError(err error, status int) *AppError {
	return &AppError{
		Error:  ErrorJSON{Err: err},
		Status: status,
	}
}

package types

import "encoding/json"

type Source struct {
	Pointer string `json:"pointer"`
}

type ErrorDetail struct {
	Source `json:"source"`
	Detail string `json:"detail"`
}

func NewValidationError(attr, message string) ErrorDetail {
	return ErrorDetail{
		Source: Source{Pointer: "data/attributes/" + attr},
		Detail: message,
	}
}

type ValidationError []ErrorDetail

func (v ValidationError) Error() string {
	errs, err := json.Marshal(struct {
		ValidationError `json:"errors"`
	}{v})

	if err != nil {
		return err.Error()
	}

	return string(errs)
}

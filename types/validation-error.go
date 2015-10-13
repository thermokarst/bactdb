package types

import "encoding/json"

type ValidationError map[string][]string

func (v ValidationError) Error() string {
	errs, err := json.Marshal(struct {
		ValidationError `json:"errors"`
	}{v})

	if err != nil {
		return err.Error()
	}

	return string(errs)
}

package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// NullFloat64 wraps sql.NullBool so that the JSON serialization can be overridden.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON makes NullFloat64 a json.Marshaller.
func (f *NullFloat64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

// UnmarshalJSON makes NullFloat64 a json.Unmarshaller.
func (f *NullFloat64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		f.Float64 = 0
		f.Valid = false
		return nil
	}
	var x interface{}
	var err error
	json.Unmarshal(b, &x)
	switch x.(type) {
	case float64:
		err = json.Unmarshal(b, &f.Float64)
	case map[string]interface{}:
		err = json.Unmarshal(b, &f.NullFloat64)
	}
	f.Valid = true
	return err
}

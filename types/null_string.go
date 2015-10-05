package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// NullString wraps sql.NullString so that the JSON serialization can be overridden.
type NullString struct {
	sql.NullString
}

// MarshalJSON makes NullString a json.Marshaller.
func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

// UnmarshalJSON makes NullString a json.Unmarshaller.
func (s *NullString) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		s.String = ""
		s.Valid = false
		return nil
	}
	var x interface{}
	var err error
	json.Unmarshal(b, &x)
	switch x.(type) {
	case string:
		err = json.Unmarshal(b, &s.String)
	case map[string]interface{}:
		err = json.Unmarshal(b, &s.NullString)
	}
	s.Valid = true
	return err
}

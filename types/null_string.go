package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

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

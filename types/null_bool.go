package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

// NullBool wraps sql.NullBool so that the JSON serialization can be overridden.
type NullBool struct {
	sql.NullBool
}

// MarshalJSON makes NullBool a json.Marshaller.
func (n *NullBool) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Bool)
}

// UnmarshalJSON makes NullBool a json.Unmarshaller.
func (n *NullBool) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		n.Bool = false
		n.Valid = false
		return nil
	}
	var x interface{}
	var err error
	json.Unmarshal(b, &x)
	switch x.(type) {
	case bool:
		err = json.Unmarshal(b, &n.Bool)
	case map[string]interface{}:
		err = json.Unmarshal(b, &n.NullBool)
	}
	n.Valid = true
	return err
}

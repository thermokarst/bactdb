package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (b *NullBool) MarshalJSON() ([]byte, error) {
	if !b.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(b.Bool)
}

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

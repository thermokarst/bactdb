package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
)

type NullInt64 struct {
	sql.NullInt64
}

func (i *NullInt64) MarshalJSON() ([]byte, error) {
	if !i.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(i.Int64)
}

func (i *NullInt64) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		i.Int64 = 0
		i.Valid = false
		return nil
	}
	var x interface{}
	var err error
	json.Unmarshal(b, &x)
	switch x.(type) {
	case float64:
		err = json.Unmarshal(b, &i.Int64)
	case map[string]interface{}:
		err = json.Unmarshal(b, &i.NullInt64)
	}
	i.Valid = true
	return err
}

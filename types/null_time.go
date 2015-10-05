package types

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
)

// NullTime wraps pq.NullTime so that the JSON serialization can be overridden.
type NullTime struct {
	pq.NullTime
}

// MarshalJSON makes NullTime a json.Marshaller.
func (t *NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

// UnmarshalJSON makes NullTime a json.Unmarshaller.
func (t *NullTime) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		var nt time.Time
		t.Time = nt.In(time.UTC)
		t.Valid = false
		return nil
	}
	var x interface{}
	var err error
	json.Unmarshal(b, &x)
	switch x.(type) {
	case string:
		err = json.Unmarshal(b, &t.Time)
	}
	t.Valid = true
	return err
}

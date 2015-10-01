package types

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
)

type NullTime struct {
	pq.NullTime
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(t.Time)
}

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

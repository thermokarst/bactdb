package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/lib/pq"
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

type NullFloat64 struct {
	sql.NullFloat64
}

func (f *NullFloat64) MarshalJSON() ([]byte, error) {
	if !f.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(f.Float64)
}

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
	case time.Time:
		err = json.Unmarshal(b, &t.Time)
	case map[string]interface{}:
		err = json.Unmarshal(b, &t.NullTime)
	}
	t.Valid = true
	return err
}

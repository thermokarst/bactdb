package types

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
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
	case string:
		err = json.Unmarshal(b, &t.Time)
	}
	t.Valid = true
	return err
}

type NullSliceInt64 []int64

func (i *NullSliceInt64) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}
	asString := string(asBytes)
	(*i) = strToIntSlice(asString)
	return nil
}

func strToIntSlice(s string) []int64 {
	r := strings.Trim(s, "{}")
	a := []int64(nil)
	for _, t := range strings.Split(r, ",") {
		if t != "NULL" {
			i, _ := strconv.ParseInt(t, 10, 64)
			a = append(a, i)
		}
	}
	return a
}

type ErrorJSON struct {
	Err error
}

func (ej ErrorJSON) Error() string {
	e, _ := json.Marshal(struct {
		Err string `json:"error"`
	}{
		Err: ej.Err.Error(),
	})
	return string(e)
}

type AppError struct {
	Error  error
	Status int
}

func NewJSONError(err error, status int) *AppError {
	return &AppError{
		Error:  ErrorJSON{Err: err},
		Status: status,
	}
}

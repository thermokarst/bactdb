package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lib/pq"
)

var (
	ErrMustProvideOptions     = errors.New("Must provide necessary options")
	ErrMustProvideOptionsJSON = newJSONError(ErrMustProvideOptions, http.StatusBadRequest)
)

// ListOptions specifies general pagination options for fetching a list of results
type ListOptions struct {
	PerPage int64   `url:",omitempty" json:",omitempty"`
	Page    int64   `url:",omitempty" json:",omitempty"`
	Ids     []int64 `url:",omitempty" json:",omitempty" schema:"ids[]"`
	Genus   string
}

func (o ListOptions) PageOrDefault() int64 {
	if o.Page <= 0 {
		return 1
	}
	return o.Page
}

func (o ListOptions) Offset() int64 {
	return (o.PageOrDefault() - 1) * o.PerPageOrDefault()
}

func (o ListOptions) PerPageOrDefault() int64 {
	if o.PerPage <= 0 {
		return DefaultPerPage
	}
	return o.PerPage
}

// DefaultPerPage is the default number of items to return in a paginated result set
const DefaultPerPage = 10

func valsIn(attribute string, values []int64, vals *[]interface{}, counter *int64) string {
	if len(values) == 1 {
		return fmt.Sprintf("%v=%v", attribute, values[0])
	}

	m := fmt.Sprintf("%v IN (", attribute)
	for _, id := range values {
		m = m + fmt.Sprintf("$%v,", *counter)
		*vals = append(*vals, id)
		*counter++
	}
	m = m[:len(m)-1] + ")"
	return m
}

func currentTime() NullTime {
	return NullTime{
		pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

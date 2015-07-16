package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/lib/pq"
)

var (
	ErrMustProvideOptions     = errors.New("Must provide necessary options")
	ErrMustProvideOptionsJSON = newJSONError(ErrMustProvideOptions, http.StatusBadRequest)
	StatusUnprocessableEntity = 422
	MustProvideAValue         = "Must provide a value"
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

func generateNonce() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func getClaims(r *http.Request) Claims {
	con := context.Get(r, "claims")
	var claims Claims
	if con != nil {
		claims = con.(Claims)
	}
	origin := r.Header.Get("Origin")
	if origin != "" {
		claims.Ref = origin
	}
	return claims
}

func canAdd(claims *Claims) bool {
	return claims.Role == "A" || claims.Role == "W"
}

func canEdit(claims *Claims, author int64) bool {
	return claims.Sub == author || claims.Role == "A"
}

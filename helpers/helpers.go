package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/schema"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/thermokarst/bactdb/types"
)

var (
	ErrMustProvideOptions     = errors.New("Must provide necessary options")
	ErrMustProvideOptionsJSON = types.NewJSONError(ErrMustProvideOptions, http.StatusBadRequest)
	StatusUnprocessableEntity = 422
	MustProvideAValue         = "Must provide a value"
	SchemaDecoder             = schema.NewDecoder()
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

type MeasurementListOptions struct {
	ListOptions
	Strains         []int64 `schema:"strain_ids"`
	Characteristics []int64 `schema:"characteristic_ids"`
}

// DefaultPerPage is the default number of items to return in a paginated result set
const DefaultPerPage = 10

func ValsIn(attribute string, values []int64, vals *[]interface{}, counter *int64) string {
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

func CurrentTime() types.NullTime {
	return types.NullTime{
		pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

// TODO: move this
func GenerateNonce() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetClaims(r *http.Request) types.Claims {
	con := context.Get(r, "claims")
	var claims types.Claims
	if con != nil {
		claims = con.(types.Claims)
	}
	origin := r.Header.Get("Origin")
	if origin != "" {
		claims.Ref = origin
	}
	return claims
}

func CanAdd(claims *types.Claims) bool {
	return claims.Role == "A" || claims.Role == "W"
}

func CanEdit(claims *types.Claims, author int64) bool {
	return claims.Sub == author || claims.Role == "A"
}

package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/context"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/gorilla/schema"
	"github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/lib/pq"
	"github.com/thermokarst/bactdb/types"
)

var (
	// StatusUnprocessableEntity is the HTTP status when Unprocessable Entity.
	StatusUnprocessableEntity = 422
	// MustProvideAValue when value required.
	MustProvideAValue = "Must provide a value"
	// SchemaDecoder for decoding schemas.
	SchemaDecoder = schema.NewDecoder()
)

// ListOptions specifies general pagination options for fetching a list of results
type ListOptions struct {
	PerPage int64   `url:",omitempty" json:",omitempty"`
	Page    int64   `url:",omitempty" json:",omitempty"`
	IDs     []int64 `url:",omitempty" json:",omitempty" schema:"ids[]"`
	Genus   string
}

// MeasurementListOptions is an extension of ListOptions.
type MeasurementListOptions struct {
	ListOptions
	Strains         []int64 `schema:"strain_ids"`
	Characteristics []int64 `schema:"characteristic_ids"`
}

// ValsIn emits X IN (A, B, C) SQL statements
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

// CurrentTime returns current time
func CurrentTime() types.NullTime {
	return types.NullTime{
		pq.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
}

// GenerateNonce generates a nonce
func GenerateNonce() (string, error) {
	// TODO: move this
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetClaims gets request claims from Authorization header
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

// CanAdd is an authorization helper for adding new entities
func CanAdd(claims *types.Claims) bool {
	return claims.Role == "A" || claims.Role == "W"
}

// CanEdit is an authorization helper for editing entities
func CanEdit(claims *types.Claims, author int64) bool {
	return claims.Sub == author || claims.Role == "A"
}

package types

import (
	"strconv"
	"strings"

	"github.com/thermokarst/bactdb/errors"
)

// NullSliceInt64 allows bactdb to read Postgres array types.
type NullSliceInt64 []int64

// Scan makes NullSliceInt64 a sql.Scanner.
func (i *NullSliceInt64) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.ErrSourceNotByteSlice
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

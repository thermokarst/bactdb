package types

import (
	"strconv"
	"strings"

	"github.com/thermokarst/bactdb/errors"
)

type NullSliceInt64 []int64

func (i *NullSliceInt64) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.SourceNotByteSlice
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

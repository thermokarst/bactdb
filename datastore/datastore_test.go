package datastore

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

func normalizeTime(t ...interface{}) {
	for _, v := range t {
		switch u := v.(type) {
		default:
			fmt.Printf("unexpected type %T", u)
		case *time.Time:
			x, _ := v.(*time.Time)
			*x = x.In(time.UTC)
		case *pq.NullTime:
			x, _ := v.(*pq.NullTime)
			*x = pq.NullTime{Time: x.Time.In(time.UTC), Valid: x.Valid}
		}
	}
}

package datastore

import (
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/thermokarst/bactdb/models"
)

func normalizeTime(t ...interface{}) {
	for _, v := range t {
		switch u := v.(type) {
		default:
			fmt.Printf("unexpected type %T", u)
		case *time.Time:
			x, _ := v.(*time.Time)
			*x = x.In(time.UTC)
		case *models.NullTime:
			x, _ := v.(*models.NullTime)
			*x = models.NullTime{pq.NullTime{Time: x.Time.In(time.UTC), Valid: x.Valid}}
		}
	}
}

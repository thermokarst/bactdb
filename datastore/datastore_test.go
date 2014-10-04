package datastore

import "time"

func normalizeTime(t ...*time.Time) {
	for _, v := range t {
		*v = v.In(time.UTC)
	}
}

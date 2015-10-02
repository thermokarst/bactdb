package models

import "github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"

var (
	// DB is a sqlx/modl database map.
	DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}
	// DBH is a global database handler.
	DBH modl.SqlExecutor = DB
)

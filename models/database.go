package models

import "github.com/thermokarst/bactdb/Godeps/_workspace/src/github.com/jmoiron/modl"

var (
	DB                   = &modl.DbMap{Dialect: modl.PostgresDialect{}}
	DBH modl.SqlExecutor = DB
)

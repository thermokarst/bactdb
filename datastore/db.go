package datastore

import (
	"log"
	"os"
	"sync"

	"github.com/DavidHuie/gomigrate"
	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DB is the global database
var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}

// DBH is a modl.SqlExecutor interface to DB, the global database. It is better to
// use DBH instead of DB because it prevents you from calling methods that could
// not later be wrapped in a transaction.
var DBH modl.SqlExecutor = DB

var connectOnce sync.Once

// Connect connects to the PostgreSQL database specified by the PG* environment
// variables. It calls log.Fatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		var err error
		conn := setDBCredentials()
		DB.Dbx, err = sqlx.Open("postgres", conn)
		if err != nil {
			log.Fatal("Error connecting to PostgreSQL database (using PG* environment variables): ", err)
		}
		DB.TraceOn("[modl]", log.New(os.Stdout, "bactdb:", log.Lmicroseconds))
		DB.Db = DB.Dbx.DB
	})
}

// Create the database schema. It calls log.Fatal if it encounters an error.
func Create(path string) {
	migrator, err := gomigrate.NewMigrator(DB.Dbx.DB, gomigrate.Postgres{}, path)
	if err != nil {
		log.Fatal("Error initializing migrations: ", err)
	}
	err = migrator.Migrate()
	if err != nil {
		log.Fatal("Error applying migrations: ", err)
	}
}

// Drop the database schema
func Drop(path string) {
	migrator, err := gomigrate.NewMigrator(DB.Dbx.DB, gomigrate.Postgres{}, path)
	if err != nil {
		log.Fatal("Error initializing migrations: ", err)
	}

	err = migrator.RollbackAll()
	if err != nil && err != gomigrate.NoActiveMigrations {
		log.Fatal("Error rolling back migrations: ", err)
	}
}

// transact calls fn in a DB transaction. If dbh is a transaction, then it just calls
// the function. Otherwise, it begins a transaction, rolling back on failure and
// committing on success.
func transact(dbh modl.SqlExecutor, fn func(fbh modl.SqlExecutor) error) error {
	var sharedTx bool
	tx, sharedTx := dbh.(*modl.Transaction)
	if !sharedTx {
		var err error
		tx, err = dbh.(*modl.DbMap).Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	}

	if err := fn(tx); err != nil {
		return err
	}

	if !sharedTx {
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func setDBCredentials() string {
	connection := "timezone=UTC "
	if heroku := os.Getenv("HEROKU"); heroku == "true" {
		url := os.Getenv("DATABASE_URL")
		conn, _ := pq.ParseURL(url)
		connection += conn
		connection += " sslmode=require"
	} else {
		connection += " sslmode=disable"
	}
	return connection
}

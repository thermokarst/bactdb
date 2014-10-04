package datastore

import (
	"log"
	"os"
	"strings"
)

func init() {
	// Make sure we don't run the tests on the main DB (will destroy the data)
	dbname := os.Getenv("PGDATABASE")
	if dbname == "" {
		dbname = "bactdbtest"
	}
	if !strings.HasSuffix(dbname, "test") {
		dbname += "test"
	}
	if err := os.Setenv("PGDATABASE", dbname); err != nil {
		log.Fatal(err)
	}

	// Reset DB
	Connect()
	Drop()
	Create()
}

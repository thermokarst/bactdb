package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/DavidHuie/gomigrate"
	"github.com/codegangsta/cli"
	"github.com/gorilla/schema"
	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/mailgun/mailgun-go"
)

var (
	DB                             = &modl.DbMap{Dialect: modl.PostgresDialect{}}
	DBH           modl.SqlExecutor = DB
	schemaDecoder                  = schema.NewDecoder()
	mgAccts                        = make(map[string]mailgun.Mailgun)
)

func main() {
	var connectOnce sync.Once
	connectOnce.Do(func() {
		var err error
		connection := "timezone=UTC "
		if heroku := os.Getenv("HEROKU"); heroku == "true" {
			url := os.Getenv("DATABASE_URL")
			conn, _ := pq.ParseURL(url)
			connection += conn
			connection += " sslmode=require"
		} else {
			connection += " sslmode=disable"
		}
		DB.Dbx, err = sqlx.Open("postgres", connection)
		if err != nil {
			log.Fatal("Error connecting to PostgreSQL database (using PG* environment variables): ", err)
		}
		DB.TraceOn("[modl]", log.New(os.Stdout, "bactdb:", log.Lmicroseconds))
		DB.Db = DB.Dbx.DB
	})

	app := cli.NewApp()
	app.Name = "bactdb"
	app.Usage = "a database for bacteria"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Matthew Ryan Dillon",
			Email: "mrdillon@alaska.edu",
		},
	}
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Start web server",
			Action:    cmdServe,
		},
		{
			Name:      "migrate",
			ShortName: "m",
			Usage:     "Migrate the database schema",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "drop",
					Usage: "Drop DB before migrating",
				},
				cli.StringFlag{
					Name:  "migration_path",
					Usage: "Path to migration DDL",
					Value: "./migrations",
				},
			},
			Action: cmdMigrateDb,
		},
	}
	app.Run(os.Args)
}

func cmdServe(c *cli.Context) {
	var err error

	// Set up Mailgun handlers:
	// [{"ref":"hymenobacter","domain":"hymenobacter.info","public":"abc","private":"123"}]
	type account struct {
		Ref     string
		Domain  string
		Public  string
		Private string
	}
	var accounts []account
	json.Unmarshal([]byte(os.Getenv("ACCOUNT_KEYS")), &accounts)
	log.Printf("Mailgun: %+v", accounts)

	for _, a := range accounts {
		mgAccts[a.Ref] = mailgun.NewMailgun(a.Domain, a.Private, a.Public)
	}

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8901"
	}
	httpAddr := fmt.Sprintf(":%v", addr)

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", Handler()))

	log.Print("Listening on ", httpAddr)
	err = http.ListenAndServe(httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func cmdMigrateDb(c *cli.Context) {
	migrationsPath := c.String("migration_path")
	migrator, err := gomigrate.NewMigrator(DB.Dbx.DB, gomigrate.Postgres{}, migrationsPath)
	if err != nil {
		log.Fatal("Error initializing migrations: ", err)
	}

	if c.Bool("drop") {
		if err = migrator.RollbackAll(); err != nil && err != gomigrate.NoActiveMigrations {
			log.Fatal("Error rolling back migrations: ", err)
		}
	}
	// Run migrations
	if err = migrator.Migrate(); err != nil {
		log.Fatal("Error applying migrations: ", err)
	}
}

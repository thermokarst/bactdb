package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/thermokarst/bactdb/api"
	"github.com/thermokarst/bactdb/datastore"
)

func main() {
	app := cli.NewApp()
	app.Name = "bactdb"
	app.Usage = "a database for bacteria"

	app.Commands = []cli.Command{
		{
			Name:      "serve",
			ShortName: "s",
			Usage:     "Start web server",
			Action:    cmdServe,
		},
		{
			Name:      "createdb",
			ShortName: "c",
			Usage:     "create the database schema",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "drop",
					Usage: "drop DB before creating",
				},
				cli.StringFlag{
					Name:  "migration_path",
					Usage: "path to migrations",
					Value: "./datastore/migrations",
				},
			},
			Action: cmdCreateDB,
		},
	}

	app.Run(os.Args)
}

func cmdServe(c *cli.Context) {
	var err error

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = "8901"
	}
	httpAddr := fmt.Sprintf(":%v", addr)

	datastore.Connect()
	err = api.SetupCerts()
	if err != nil {
		log.Fatal("SetupCerts: ", err)
	}

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", corsHandler(api.Handler())))
	log.Print("Listening on ", httpAddr)
	err = http.ListenAndServe(httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func cmdCreateDB(c *cli.Context) {
	migrationsPath := c.String("migration_path")

	datastore.Connect()

	if c.Bool("drop") {
		datastore.Drop(migrationsPath)
	}
	datastore.Create(migrationsPath)
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domains := os.Getenv("DOMAINS")
		allowedDomains := strings.Split(domains, ",")
		if origin := r.Header.Get("Origin"); origin != "" {
			for _, s := range allowedDomains {
				if s == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
					w.Header().Set("Access-Control-Allow-Methods", r.Header.Get("Access-Control-Request-Method"))
				}
			}
		}
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	}
}

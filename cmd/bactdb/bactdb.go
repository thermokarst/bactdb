package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Usage: "HTTP service port",
					Value: 8901,
				},
				cli.StringFlag{
					Name:  "keys",
					Usage: "path to keys",
					Value: "keys/",
				},
			},
			Action: cmdServe,
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
	httpAddr := fmt.Sprintf(":%v", c.Int("port"))

	datastore.Connect()
	api.SetupCerts(c.String("keys"))

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	log.Print("Listening on ", httpAddr)
	err := http.ListenAndServe(httpAddr, m)
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

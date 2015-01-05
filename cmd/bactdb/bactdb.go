package main

import (
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
				cli.StringFlag{
					Name:  "http",
					Value: ":8901",
					Usage: "HTTP service address",
				},
			},
			Action: cmdServe,
		},
	}

	app.Run(os.Args)
}

func cmdServe(c *cli.Context) {
	httpAddr := c.String("http")

	datastore.Connect()

	m := http.NewServeMux()
	m.Handle("/api/", http.StripPrefix("/api", api.Handler()))
	log.Print("Listening on ", httpAddr)
	err := http.ListenAndServe(httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

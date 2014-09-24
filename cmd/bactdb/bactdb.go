package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thermokarst/bactdb/api"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `bactdb is a database for bacteria.

Usage:

		bactdb [options] command [arg...]

The commands are:
`)
		for _, c := range subcmds {
			fmt.Fprintf(os.Stderr, "    %-24s %s\n", c.name, c.description)
		}
		fmt.Fprintln(os.Stderr, `
Use "bactdb command -h" for more information about a command.

The options are:
`)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}
	log.SetFlags(0)

	subcmd := flag.Arg(0)
	for _, c := range subcmds {
		if c.name == subcmd {
			c.run(flag.Args()[1:])
			return
		}
	}

	fmt.Fprintf(os.Stderr, "unknown subcmd %q\n", subcmd)
	fmt.Fprintln(os.Stderr, `Run "bactdb -h" for usage.`)
	os.Exit(1)
}

type subcmd struct {
	name        string
	description string
	run         func(args []string)
}

var subcmds = []subcmd{
	{"serve", "start web server", serveCmd},
}

func serveCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	httpAddr := flag.String("http", ":8901", "HTTP service address")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bactdb serve [options]

Starts the web server that serves the API.

The options are:
`)
		fs.PrintDefaults()
		os.Exit(1)
	}
	fs.Parse(args)

	if fs.NArg() != 0 {
		fs.Usage()
	}

	m := http.NewServeMux()
	m.Handle("/api", api.Handler())

	log.Print("Listening on ", *httpAddr)
	err := http.ListenAndServe(*httpAddr, m)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/johnernaut/webhog/webhog/router"
	"github.com/martini-contrib/binding"
	_ "github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"net/http"
	"os"
	"runtime"
)

const VERSION = "v0.1.0"

func main() {
	version := flag.Bool("version", false, "current version")
	flag.Parse()
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// All the parallelism are belong to us!
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load configuration file.
	webhog.LoadConfig()

	// Load DB instance.
	webhog.LoadDB()

	// Start the server.
	m := martini.Classic()

	m.Use(render.Renderer())
	// m.Use(gzip.All())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))

	m.Group("/api", func(r martini.Router) {
		r.Post("/scrape", KeyRequired(), binding.Bind(router.Url{}), router.Scrape)

		r.Get("/entity/:uuid", router.Entity)

		r.Get("/entities", router.Entities)
	})

	m.Run()
}

func KeyRequired() martini.Handler {
	return func(context martini.Context, res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-KEY") != webhog.Config.ApiKey {
			http.Error(res, "Invalid API key.", http.StatusForbidden)
		}
		context.Next()
	}
}

package main

import (
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"runtime"
)

type Url struct {
	Url  string `json:"url"`
	UUID string `json:"uuid"`
}

func main() {
	// All the parallelism are belong to us!
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Load configuration file.
	webhog.LoadConfig()

	// Start the server.
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Use(webhog.DB())

	// Middleware to make sure each request has a specified API key
	m.Use(func(res http.ResponseWriter, req *http.Request, r render.Render) {
		if req.Header.Get("X-API-KEY") != webhog.Config.ApiKey {
			r.JSON(401, map[string]interface{}{"error": "Invalid API key."})
		}
	})

	m.Post("/scrape", binding.Json(Url{}), func(url Url, r render.Render) {
		entity, err := webhog.NewScraper(url.Url)
		if err != nil {
			r.JSON(400, map[string]interface{}{"errors": err.Error()})
		} else {
			r.JSON(200, entity)
		}
	})

	m.Run()
}

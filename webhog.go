package main

import (
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/martini-contrib/binding"
	_ "github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo/bson"
	"log"
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

	// Load DB instance.
	webhog.LoadDB()

	// Start the server.
	m := martini.Classic()

	m.Use(render.Renderer())
	// m.Use(gzip.All())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))

	m.Group("/api", func(r martini.Router) {
		r.Post("/scrape", binding.Json(Url{}), func(url Url, r render.Render) {
			log.Println("SUP")
			entity, err := webhog.NewScraper(url.Url)
			log.Println(entity)
			if err != nil {
				r.JSON(400, map[string]interface{}{"errors": err.Error()})
			} else {
				r.JSON(200, entity)
			}
		})

		r.Get("/entity/:uuid", func(params martini.Params, r render.Render) {
			entity := new(webhog.Entity)
			err := entity.Find(bson.M{"uuid": params["uuid"]})

			if err != nil {
				r.JSON(200, map[string]interface{}{"errors": "Entity not found."})
			} else {
				r.JSON(200, entity)
			}
		})
	})

	m.Run()
}

func KeyRequired() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, r render.Render) {
		if req.Header.Get("X-API-KEY") != webhog.Config.ApiKey {
			r.JSON(401, map[string]interface{}{"error": "Invalid API key."})
		}
	}
}

package main

import (
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/martini-contrib/binding"
	_ "github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"runtime"
)

type Url struct {
	Url  string `form:"url" json:"url"`
	UUID string `form:"uuid" json:"uuid"`
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
		r.Post("/scrape", binding.Bind(Url{}), func(url Url, r render.Render) {
			entity, err := webhog.NewScraper(url.Url)
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

func (urlType Url) Validate(errors binding.Errors, req *http.Request) binding.Errors {
	_, err := url.ParseRequestURI(urlType.Url)
	if err != nil {
		errors = append(errors, binding.Error{
			Message: "Malformed URL. Please provide proper URL formatting.",
		})
	}
	return errors
}

package router

import (
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
)

func LoadRoutes() {
	// Start the server.
	m := martini.Classic()

	m.Use(render.Renderer())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))

	m.Group("/api", func(r martini.Router) {
		r.Post("/scrape", KeyRequired(), binding.Bind(Url{}), Scrape)

		r.Get("/entity/:uuid", Entity)

		r.Get("/entities", Entities)
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

type Url struct {
	Url  string `form:"url" json:"url"`
	UUID string `form:"uuid" json:"uuid"`
}

func Scrape(url Url, r render.Render) {
	entity, err := webhog.NewScraper(url.Url)
	if err != nil {
		r.JSON(400, map[string]interface{}{"errors": err.Error()})
	} else {
		r.JSON(200, entity)
	}
}

func Entity(params martini.Params, r render.Render) {
	entity := new(webhog.Entity)
	err := entity.Find(bson.M{"uuid": params["uuid"]})

	if err != nil {
		r.JSON(200, map[string]interface{}{"errors": "Entity not found."})
	} else {
		r.JSON(200, entity)
	}
}

func Entities(params martini.Params, r render.Render) {
	entity := new(webhog.Entity)
	entities, err := entity.All()

	if err != nil {
		r.JSON(200, map[string]interface{}{"errors": "Entity not found."})
	} else {
		r.JSON(200, entities)
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

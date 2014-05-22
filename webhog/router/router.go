package router

import (
	"github.com/go-martini/martini"
	"github.com/johnernaut/webhog/webhog"
	"github.com/martini-contrib/binding"
	_ "github.com/martini-contrib/gzip"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
)

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

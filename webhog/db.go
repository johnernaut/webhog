package webhog

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
)

// Interface that wraps DB models for a common
// querying interface.
type Model interface {
	Find(interface{}, *mgo.Collection) error
	Create(*mgo.Collection) error
	Update(interface{}, interface{}, *mgo.Collection) error
}

// Hold a reference to all models.
var Models = []Model{}

// Register a model object into the Models
// reference.
func Register(m Model) {
	Models = append(Models, m)
}

// Connect to the given database and pass the connection into
// martini
func DB() martini.Handler {
	session, err := mgo.Dial(Config.mongodb)
	if err != nil {
		log.Panicln("Error establishing database connection: ", err)
	}

	return func(c martini.Context) {
		s := session.Clone()

		c.Map(s.DB("webhog").C("entities"))
		defer s.Close()
		c.Next()
	}
}

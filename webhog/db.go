package webhog

import (
	"github.com/go-martini/martini"
	"labix.org/v2/mgo"
	"log"
)

// Type to hold DB connection information
type connection struct {
	Sess *mgo.Session
	C    *mgo.Collection
}

var Db *connection

// Interface that wraps DB models for a common
// querying interface.
type Model interface {
	Find(interface{}) error
	Create() error
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
	Db = new(connection)
	session, err := mgo.Dial(Config.mongodb)
	if err != nil {
		log.Panicln("Error establishing database connection: ", err)
	}

	return func(c martini.Context) {
		s := session.Clone()

		Db.Sess = s
		Db.C = s.DB("webhog").C("entities")

		c.Map(Db.C)
		defer s.Close()
		c.Next()
	}
}

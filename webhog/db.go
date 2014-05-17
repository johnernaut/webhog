package webhog

import (
	"labix.org/v2/mgo"
	"log"
)

// Type to hold Mongo DB connection info
type connection struct {
	Sess *mgo.Session
	C    *mgo.Collection
}

// Global var to hold the DB connection
var Conn = new(connection)

// Interface that wraps DB models for a common
// querying interface.
type Model interface {
	Find(interface{}) error
	Create() error
	Update(interface{}, interface{}) error
	All() ([]Entity, error)
}

// Hold a reference to all models.
var Models = []Model{}

// Register a model object into the Models
// reference.
func Register(m Model) {
	Models = append(Models, m)
}

// Connect to the given database
func LoadDB() {
	session, err := mgo.Dial(Config.mongodb)
	if err != nil {
		log.Panicln("Error establishing database connection: ", err)
	}

	Conn.Sess = session
	Conn.C = session.DB("webhog").C("entities")
}

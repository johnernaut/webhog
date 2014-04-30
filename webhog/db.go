package webhog

import (
	"labix.org/v2/mgo"
	"log"
)

// Holds an instance of the MySQL DB connection.
type connection struct {
	Db *mgo.Database
}

// Interface that wraps DB models for a common
// querying interface.
type Model interface {
	Find(string) *Entity
	Create() *Entity
}

// Global var that references a handler to the
// DB connection.
var Connection = new(connection)

// Hold a reference to all models.
var Models = []Model{}

// Register a model object into the Models
// reference.
func Register(m Model) {
	Models = append(Models, m)
}

// Connect to the given database and pass the connection into
// our Connection object for re-use.
func ConnectDB() {
	session, err := mgo.Dial(Config.mongodb)
	if err != nil {
		log.Panic(err)
	}

	defer session.Close()

	Connection.Db = session.DB("webhog")
}

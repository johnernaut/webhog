package webhog

import (
	"labix.org/v2/mgo"
	"log"
)

// Type to hold Mongo DB connection info
type connection struct {
	Db *mgo.Database
}

// Global var to hold the DB connection
var Conn = new(connection)

// Interface that wraps DB models for a common
// querying interface.
type Model interface {
	Collection() string
}

// Hold a reference to all models.
var Models = []Model{}

// Register a model object into the Models
// reference.
func Register(m Model) {
	Models = append(Models, m)
}

func Cursor(m Model) *mgo.Collection {
	return Conn.Db.C(m.Collection())
}

func Find(m Model, query interface{}) *mgo.Query {
	cursor := Cursor(m)
	return cursor.Find(query)
}

func Update(m Model, query, updates interface{}) error {
	cursor := Cursor(m)
	err := cursor.Update(query, updates)
	return err
}

func Create(m Model) error {
	cursor := Cursor(m)
	err := cursor.Insert(m)
	if err != nil {
		return err
	}

	return err
}

// Connect to the given database
func LoadDB() {
	session, err := mgo.Dial(Config.mongodb)
	if err != nil {
		log.Panicln("Error establishing database connection: ", err)
	}

	Conn.Db = session.DB("")
}

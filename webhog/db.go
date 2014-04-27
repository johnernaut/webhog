package webhog

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// Holds an instance of the MySQL DB connection.
type connection struct {
	Db *sql.DB
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
	tcpString := "@tcp(" + Config.mysql + ")"
	db, err := sql.Open("mysql", "root:"+tcpString+"/"+Config.dbName)

	if err != nil {
		log.Panic(err)
	}

	Connection.Db = db

	err = Connection.Db.Ping()

	if err != nil {
		log.Panic(err)
	}
}

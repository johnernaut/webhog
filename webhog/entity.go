package webhog

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

// Entity is a representation of a webpage and it's corresponding
// UUID that's stored on AWS-S3
type Entity struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UUID      string        `bson:"uuid" json:"uuid"`
	Url       string        `bson:"url" json:"url"`
	AwsLink   string        `bson:"aws_link,omitempty" json:"aws_link"`
	Status    string        `bson:"status" json:"status"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

// Find an Entity by UUID and return the AWS S3 link (or status if upload
// is incomplete) - create if it doesn't exist.
func (entity *Entity) Find(query string, db *mgo.Collection) (result *Entity) {
	// query is url
	err := db.Find(bson.M{"url": query}).One(&entity)
	if err != nil {
		return entity
	}

	return entity
}

// Create a new entity object in the database.
func (entity *Entity) Create(db *mgo.Collection) (result *Entity) {
	err := db.Insert(entity)
	if err != nil {
		log.Panic(err)
	}

	return entity
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

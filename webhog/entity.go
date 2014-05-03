package webhog

import (
	"labix.org/v2/mgo/bson"
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
func (entity *Entity) Find(query interface{}) error {
	// query is url
	err := Db.C.Find(query).One(&entity)

	return err
}

// Create a new entity object in the database.
func (entity *Entity) Create() error {
	err := Db.C.Insert(entity)

	return err
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

package webhog

import (
	"labix.org/v2/mgo/bson"
	_ "log"
	_ "time"
)

// Entity is a representation of a webpage and it's corresponding
// UUID that's stored on AWS-S3
type Entity struct {
	Id      bson.ObjectId `bson:"_id" json:"id"`
	UUID    string        `bson:"uuid" json:"uuid"`
	Url     string        `bson:"url" json:"url"`
	AwsLink string        `bson:"aws_link,omitempty" json:"aws_link"`
	Status  string        `bson:"status" json:"status"`
}

// Find an Entity by UUID and return the AWS S3 link (or status if upload
// is incomplete) - create if it doesn't exist.
func (entity *Entity) Find(query string) (result *Entity) {
	// query is url
	return nil
}

// Create a new entity object in the database.
func (entity *Entity) Create() (result *Entity) {
	return nil
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

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

// Set a URL's expiration time to 1 week before it needs
// to be reprocessed.
var ExpirationTime = time.Hour * 168

func (e Entity) Collection() string {
	return "entities"
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

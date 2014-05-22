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

// Find an Entity by UUID and return the AWS S3 link (or status if upload
// is incomplete) - create if it doesn't exist.
func (entity *Entity) Find(query interface{}) error {
	// query is url
	err := Conn.C.Find(query).One(&entity)

	return err
}

// Find all entities
func (entity *Entity) All() ([]Entity, error) {
	entities := []Entity{}
	err := Conn.C.Find(nil).All(&entities)

	return entities, err
}

// Update an entities' attributes
func (entity *Entity) Update(query interface{}, updates interface{}) error {
	// query is url
	err := Conn.C.Update(query, updates)

	return err
}

// Create a new entity object in the database.
func (entity *Entity) Create() error {
	err := Conn.C.Insert(entity)
	if err != nil {
		return err
	}

	err = Conn.C.Find(bson.M{"uuid": entity.UUID}).One(&entity)

	return err
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

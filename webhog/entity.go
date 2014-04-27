package webhog

import (
	"log"
	"time"
)

// Entity is a representation of a webpage and it's corresponding
// UUID that's stored on AWS-S3
type Entity struct {
	Id        int       `json:"id"`
	UUID      string    `field:"uuid" json:"uuid"`
	Url       string    `json:"url"`
	AwsLink   string    `field:"aws_link" json:"aws_link"`
	Status    string    `json:"status"`
	CreatedAt time.Time `field:"created_at" json:"created_at"`
	UpdatedAt time.Time `field:"updated_at" json:"updated_at"`
}

// Find an Entity by UUID and return the AWS S3 link (or status if upload
// is incomplete) - create if it doesn't exist.
func (entity *Entity) Find(query string) (result *Entity) {
	rows, err := Connection.Db.Query("SELECT id,uuid,url,aws_link,status FROM entities WHERE url=? LIMIT 1", query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		if err := rows.Scan(&entity.Id, &entity.UUID, &entity.Url, &entity.AwsLink, &entity.Status); err != nil {
			log.Println(err)
		}

		return entity
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return entity
}

// Create a new entity object in the database.
func (entity *Entity) Create() (result *Entity) {
	res, err := Connection.Db.Exec("INSERT INTO entities(uuid, url, status) VALUES(?, ?, ?)", entity.UUID, entity.Url, entity.Status)
	if err != nil {
		log.Println(err)
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println("Error in lastId: ", err)
	}

	err = Connection.Db.QueryRow("SELECT id,uuid,url,status from entities where id=?", lastId).
		Scan(&entity.Id, &entity.UUID, &entity.Url, &entity.Status)
	if err != nil {
		log.Println("Error after queryrow: ", err)
	}

	return entity
}

func init() {
	// Register Entity as a Model.
	Register(&Entity{})
}

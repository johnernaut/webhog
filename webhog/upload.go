package webhog

import (
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
	"log"
	"strings"
)

func UploadEntity(dir string, entity *Entity) error {
	spl := strings.Split(dir, "/")
	endDir := spl[len(spl)-1]

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Println(err.Error())
	}

	// Open Bucket
	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(Config.bucket)
	log.Println(dir)

	b, err := ioutil.ReadFile(dir)
	if err != nil {
		log.Println("Error reading file for upload: ", err)
	}

	err = bucket.Put("/"+endDir, b, "text/plain", s3.BucketOwnerFull)
	if err != nil {
		log.Println("Error uploading file to bucket: ", err)
	}

	return nil
}

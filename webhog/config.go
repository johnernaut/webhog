package webhog

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"log"
	"os"
	"strings"
)

type configuration struct {
	mongodb   string
	ApiKey    string
	AwsKey    string
	AwsSecret string
	bucket    string
}

var Config = new(configuration)

func LoadConfig() error {
	conf, err := yaml.ReadFile("webhog.yml")
	if err != nil {
		return err
	}

	Config.mongodb, _ = conf.Get(getEnv() + ".mongodb")
	Config.ApiKey, _ = conf.Get(getEnv() + ".api_key")
	Config.AwsKey, _ = conf.Get(getEnv() + ".aws_key")
	Config.AwsSecret, _ = conf.Get(getEnv() + ".aws_secret")
	Config.bucket, _ = conf.Get(getEnv() + ".bucket")

	f, err := os.OpenFile("webhog_"+getEnv()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error writing log file: ", err)
	}

	log.SetOutput(f)

	return err
}

func getEnv() string {
	env := strings.ToLower(os.Getenv("MARTINI_ENV"))
	if env == "" || env == "development" {
		return "development"
	}

	return env
}

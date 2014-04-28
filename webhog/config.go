package webhog

import (
	"github.com/kylelemons/go-gypsy/yaml"
	"os"
)

type configuration struct {
	mysql     string
	dbName    string
	ApiKey    string
	awsKey    string
	awsSecret string
	bucket    string
}

var Config = new(configuration)

func LoadConfig() error {
	conf, err := yaml.ReadFile("webhog.yml")
	if err != nil {
		return err
	}

	Config.mysql, _ = conf.Get(getEnv() + ".mysql")
	Config.dbName, _ = conf.Get(getEnv() + ".db_name")
	Config.ApiKey, _ = conf.Get(getEnv() + ".api_key")
	Config.awsKey, _ = conf.Get(getEnv() + ".aws_key")
	Config.awsSecret, _ = conf.Get(getEnv() + ".aws_secret")
	Config.bucket, _ = conf.Get(getEnv() + ".bucket")

	return err
}

func getEnv() string {
	env := os.Getenv("GO_ENV")
	if env == "" || env == "development" {
		return "development"
	}

	return env
}

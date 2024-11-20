package configuration

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

func BuildConfiguration() (*Configuration, error) {
	if os.Getenv("APP_MODE") == "DEBUG" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error dugirn loading .env file: %v", err)

			return nil, err
		}
	}

	var c Configuration
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}

	return &c, nil
}

type Configuration struct {
	InfluxUrl    string `envconfig:"INFLUXDB_URL"`
	InfluxOrg    string `envconfig:"INFLUXDB_ORG"`
	InfluxBucket string `envconfig:"INFLUXDB_BUCKET"`
	InfluxToken  string `envconfig:"INFLUXDB2_ADMIN_TOKEN"`

	RedisUrl      string `envconfig:"REDIS_URL"`
	RedisDB       int    `envconfig:"REDIS_DB"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`

	Port int `envconfig:"PORT" default:"8080"`
}

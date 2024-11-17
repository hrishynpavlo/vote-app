package configuration

import (
	"github.com/kelseyhightower/envconfig"
)

func BuildConfiguration() (*Configuration, error) {
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
}

package main

import (
	"github.com/rcrowley/go-metrics"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"vote-app/api"
	"vote-app/configuration"
	metricsreporter "vote-app/metrics-reporter"
	"vote-app/persistance"
)

func main() {
	fx.New(
		fx.Provide(
			BuildConfiguration,
			Persistence,
			func() metrics.Registry { return metrics.DefaultRegistry },
			metricsreporter.MetricsReporter,
			api.AddVoteController,
			api.AddGin),
		fx.Invoke(
			func(*metricsreporter.Reporter) {},
			func(server *api.Server) {},
		),
	).Run()
}

func BuildConfiguration() *configuration.Configuration {
	cfg, err := configuration.BuildConfiguration()
	if err != nil {
		panic(err)
	}

	return cfg
}
func Persistence(cfg *configuration.Configuration) *persistance.RedisCache {
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB})

	return &persistance.RedisCache{
		Db: db,
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
	"vote-app/api"
	"vote-app/configuration"
	metrics_reporter "vote-app/metrics-reporter"
)

func main() {
	cfg, err := configuration.BuildConfiguration()
	if err != nil {
		panic(err)
	}

	metricsRegistry := metrics.DefaultRegistry
	go metrics_reporter.InfluxDBWithTags(metricsRegistry, time.Second*30, cfg.InfluxUrl,
		cfg.InfluxOrg, cfg.InfluxBucket, "unit", cfg.InfluxToken, BuildDefaultMetricsTags(), true)

	db := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisUrl,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB})

	r := gin.New()

	r.Use(Metrics(metricsRegistry))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})

		return
	})

	r.POST("/vote", func(c *gin.Context) {
		api.CreateVote(c, db)
	})

	r.GET("/vote", func(c *gin.Context) {
		api.GetVotes(c, db)
	})

	err = r.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func Metrics(metricsRegistry metrics.Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestsRate := metrics.GetOrRegisterCounter("requests_number", metricsRegistry)
		latencyTracker := metrics.GetOrRegisterGauge("requests_latency", metricsRegistry)

		requestsRate.Inc(1)
		now := time.Now().UTC()
		c.Next()
		latency := time.Since(now).Microseconds()
		latencyTracker.Update(latency)
	}
}
func BuildDefaultMetricsTags() map[string]string {
	tags := map[string]string{"app": "vote_api", "machine_name": "vote_api_local"}
	return tags
}

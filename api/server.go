package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/fx"
	"net/http"
	"time"
	"vote-app/configuration"
)

type Server struct {
	r *gin.Engine
}

func AddGin(lc fx.Lifecycle, metricsRegistry metrics.Registry, voteController *VoteController, configuration *configuration.Configuration) *Server {
	r := gin.New()

	r.Use(collectMetrics(metricsRegistry))

	r.GET("/", pingPong)

	r.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, metricsRegistry.GetAll())
		return
	})

	api := r.Group("/api")
	voteApi := api.Group("/vote")

	voteApi.POST("", voteController.CreateVote)

	voteApi.GET("", voteController.GetVotes)

	voteApi.GET("/:id/stats", voteController.GetVoteStats)

	voteApi.PATCH("/:id/submit/:optionId", voteController.Vote)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := r.Run(fmt.Sprintf(":%d", configuration.Port))
				if err != nil {
					panic(err)
				}
			}()

			return nil
		},
	})

	return &Server{
		r: r,
	}
}

func collectMetrics(metricsRegistry metrics.Registry) gin.HandlerFunc {
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
func pingPong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})

	return
}

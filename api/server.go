package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
	"go.uber.org/fx"
	"net/http"
	"time"
)

type Server struct {
	r *gin.Engine
}

func AddGin(lc fx.Lifecycle, metricsRegistry metrics.Registry, voteController *VoteController) *Server {
	r := gin.New()

	r.Use(collectMetrics(metricsRegistry))

	r.GET("/", pingPong)

	r.POST("/vote", voteController.CreateVote)

	r.GET("/vote", voteController.GetVotes)

	r.GET("/vote/:id/stats", voteController.GetVoteStats)

	r.PATCH("/vote/:id/submit/:optionId", voteController.Vote)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				err := r.Run(":8080")
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

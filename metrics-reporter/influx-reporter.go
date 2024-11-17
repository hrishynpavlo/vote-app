package metrics_reporter

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	_ "github.com/influxdata/line-protocol"
	"log"
	uurl "net/url"
	"time"

	client "github.com/influxdata/influxdb-client-go/v2"
	"github.com/rcrowley/go-metrics"
)

type reporter struct {
	reg          metrics.Registry
	interval     time.Duration
	align        bool
	url          uurl.URL
	organization string
	bucket       string
	measurement  string
	token        string
	tags         map[string]string
	influxClient client.Client
	client       api.WriteAPI
}

// InfluxDBWithTags starts a InfluxDB reporter which will post the metrics from the given registry at each d interval with the specified tags
func InfluxDBWithTags(r metrics.Registry, d time.Duration, url, organization, bucket, measurement, token string, tags map[string]string, align bool) {
	u, err := uurl.Parse(url)
	if err != nil {
		log.Printf("unable to parse InfluxDB url %s. err=%v", url, err)
		return
	}

	rep := &reporter{
		reg:          r,
		interval:     d,
		url:          *u,
		organization: organization,
		bucket:       bucket,
		measurement:  measurement,
		token:        token,
		tags:         tags,
		align:        align,
	}
	if err := rep.makeClient(); err != nil {
		log.Printf("unable to make InfluxDB client. err=%v", err)
		return
	}

	rep.run()
}

func (r *reporter) makeClient() (err error) {
	r.influxClient = client.NewClientWithOptions(r.url.String(), r.token, client.DefaultOptions().SetBatchSize(5000).SetFlushInterval(uint(r.interval.Milliseconds())))
	r.client = r.influxClient.WriteAPI(r.organization, r.bucket)

	return
}

func (r *reporter) run() {
	intervalTicker := time.Tick(r.interval)
	pingTicker := time.Tick(time.Second * 5)
	errorsChannel := r.client.Errors()

	go func() {
		for err := range errorsChannel {
			log.Printf("unable to send metrics to InfluxDB. err=%v", err)
		}
	}()

	for {
		select {
		case <-intervalTicker:
			r.send()
		case <-pingTicker:
			_, err := r.influxClient.Ping(context.Background())
			if err != nil {
				log.Printf("got error while sending a ping to InfluxDB, trying to recreate client. err=%v", err)

				if err = r.makeClient(); err != nil {
					log.Printf("unable to make InfluxDB client. err=%v", err)
				}
			}
		}
	}
}

func (r *reporter) send() {
	points := make([]*write.Point, 0)
	now := time.Now()
	if r.align {
		now = now.Truncate(r.interval)
	}
	r.reg.Each(func(name string, i interface{}) {

		switch metric := i.(type) {
		case metrics.Counter:
			ms := metric.Snapshot()
			p := client.NewPoint(r.measurement, r.tags, map[string]interface{}{
				fmt.Sprintf("%s.count", name): ms.Count()}, now)
			points = append(points, p)

		case metrics.Gauge:
			ms := metric.Snapshot()
			p := client.NewPoint(r.measurement, r.tags, map[string]interface{}{
				fmt.Sprintf("%s.gauge", name): ms.Value(),
			}, now)

			points = append(points, p)

		case metrics.GaugeFloat64:
			ms := metric.Snapshot()
			p := client.NewPoint(r.measurement, r.tags, map[string]interface{}{
				fmt.Sprintf("%s.gauge", name): ms.Value(),
			}, now)

			points = append(points, p)

		case metrics.Histogram:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			fields := map[string]float64{
				"count":    float64(ms.Count()),
				"max":      float64(ms.Max()),
				"mean":     ms.Mean(),
				"min":      float64(ms.Min()),
				"stddev":   ms.StdDev(),
				"variance": ms.Variance(),
				"p50":      ps[0],
				"p75":      ps[1],
				"p95":      ps[2],
				"p99":      ps[3],
				"p999":     ps[4],
				"p9999":    ps[5],
			}
			for k, v := range fields {
				p := client.NewPoint(r.measurement, bucketTags(k, r.tags), map[string]interface{}{
					fmt.Sprintf("%s.histogram", name): v,
				}, now)

				points = append(points, p)
			}

		case metrics.Meter:
			ms := metric.Snapshot()
			fields := map[string]float64{
				"count": float64(ms.Count()),
				"m1":    ms.Rate1(),
				"m5":    ms.Rate5(),
				"m15":   ms.Rate15(),
				"mean":  ms.RateMean(),
			}
			for k, v := range fields {
				p := client.NewPoint(r.measurement, bucketTags(k, r.tags), map[string]interface{}{
					fmt.Sprintf("%s.meter", name): v,
				}, now)

				points = append(points, p)
			}

		case metrics.Timer:
			ms := metric.Snapshot()
			ps := ms.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999})
			fields := map[string]float64{
				"count":    float64(ms.Count()),
				"max":      float64(ms.Max()),
				"mean":     ms.Mean(),
				"min":      float64(ms.Min()),
				"stddev":   ms.StdDev(),
				"variance": ms.Variance(),
				"p50":      ps[0],
				"p75":      ps[1],
				"p95":      ps[2],
				"p99":      ps[3],
				"p999":     ps[4],
				"p9999":    ps[5],
				"m1":       ms.Rate1(),
				"m5":       ms.Rate5(),
				"m15":      ms.Rate15(),
				"meanrate": ms.RateMean(),
			}
			for k, v := range fields {
				p := client.NewPoint(r.measurement, bucketTags(k, r.tags), map[string]interface{}{
					fmt.Sprintf("%s.timer", name): v,
				}, now)

				points = append(points, p)
			}
		}
	})

	for _, p := range points {
		r.client.WritePoint(p)
	}

	return
}

func bucketTags(bucket string, tags map[string]string) map[string]string {
	m := map[string]string{}
	for tk, tv := range tags {
		m[tk] = tv
	}
	m["bucket"] = bucket
	return m
}

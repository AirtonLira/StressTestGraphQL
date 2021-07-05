package metric

import (
	"fmt"
	"log"
	"time"

	"github.com/google/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//HTTP application
type HTTP struct {
	Query      string
	StatusCode int
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

var (
	opsRequested = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stresstest_graphql_ops_total",
		Help: "The total number of request GraphQL",
	})

	http = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "http",
		Name:      "request_duration_seconds",
		Help:      "The latency of the HTTP requests.",
	}, []string{"StartedAt", "Query", "StatusCode"})
)

//TimeTrack print time to request
func (h *HTTP) TimeTrack(start time.Time, name string, statusCode int) {
	elapsed := time.Since(start)
	h.StatusCode = statusCode
	h.Query = name
	h.StartedAt = start
	h.FinishedAt = time.Now()
	h.Duration = elapsed.Seconds()
	opsRequested.Inc()

	prometheus.Register(http)
	http.WithLabelValues(h.StartedAt.Format("02-Jan-2006"), h.Query, fmt.Sprint(h.StatusCode)).Set(h.Duration)

	log.Printf("Query: %s took %f - StatusCode: %d", name, elapsed.Seconds(), h.StatusCode)
	logger.Info("Query: ", name, " took ", elapsed.Seconds(), " - StatusCode: ", h.StatusCode)
}

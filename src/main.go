package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_app_requests_total",
		Help: "Total number of HTTP requests",
	})

	requestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "go_app_request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.1, 0.5, 1, 2, 5},
	})
)

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		requestDuration.Observe(duration)
	}()

	requestsTotal.Inc()
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", handler)

	prometheus.Register(prometheus.NewGoCollector())
	prometheus.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
aws-dc-exporter is used to fetch metrics from AWS Direct Connect APIs and export
them as Prometheus metrics. These metrics can be plugged into an alerting solution
and create alerts like, `When a BGP router is down?` or `AWS Connection State is not up`.
Usage Instructions
`./aws-dc-exporter --config=config.yml`
*/

package main

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/VictoriaMetrics/metrics"
)

var (
	// injected during build
	buildVersion = "unknown"
	buildDate    = "unknown"
)

func initLogger(config cfgApp) *logrus.Logger {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// Set logger level
	switch level := config.LogLevel; level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
		logger.Debug("verbose logging enabled")
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
	return logger
}

func main() {
	var (
		config = initConfig()
		logger = initLogger(config.App)
	)
	// Initialize hub.
	hub := &Hub{
		config:  config,
		logger:  logger,
		version: buildVersion,
	}
	hub.logger.Infof("booting aws-dc-exporter version:%v", buildVersion)
	// Fetch all jobs listed in config.
	for _, job := range hub.config.App.Jobs {
		// This is to avoid all copies of `exporter` getting updated by the last `job` memory address
		// you instantiate with, since we pass `job` as a pointer to the struct.
		j := job
		// Initialize the exporter. Exporter is a collection of metrics to be exported.
		_, err := hub.NewExporter(&j)
		if err != nil {
			hub.logger.Errorf("exporter initialization failed for %s", job.Name)
		}
	}
	// Default index handler.
	handleIndex := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to aws-dc-exporter. Visit /metrics."))
	})
	// Metrics handler
	handleMetrics := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	// Initialize router and define all endpoints.
	router := http.NewServeMux()
	router.Handle("/", handleIndex)
	// Expose the registered metrics at `/metrics` path.
	router.Handle("/", handleMetrics)

	// Initialize server.
	server := &http.Server{
		Addr:         hub.config.Server.Address,
		Handler:      router,
		ReadTimeout:  hub.config.Server.ReadTimeout * time.Millisecond,
		WriteTimeout: hub.config.Server.WriteTimeout * time.Millisecond,
	}
	// Start the server. Blocks the main thread.
	hub.logger.Infof("starting server listening on %v", hub.config.Server.Address)
	if err := server.ListenAndServe(); err != nil {
		hub.logger.Fatalf("error starting server: %v", err)
	}
}

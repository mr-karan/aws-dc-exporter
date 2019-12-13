package main

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
)

const (
	namespace = "aws_dc"
)

// NewExporter returns an initialized `Exporter`.
func (hub *Hub) NewExporter(job *Job) (*Exporter, error) {
	dc, err := hub.NewDCClient(&job.AWSCreds)
	if err != nil {
		hub.logger.Errorf("Error initializing AWS Client")
		return nil, err
	}
	return &Exporter{
		client: dc,
		job:    job,
		hub:    hub,
		up:     metrics.GetOrCreateGauge(fmt.Sprintf(`%s{job="%s"}`, namespace, job.Name), up),
	}, nil
}

func up() float64 {
	return 1
}

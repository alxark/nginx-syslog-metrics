package internal

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strconv"
	"strings"
)

type Counter struct {
	log    *log.Logger
	Prefix string
}

func NewCounter(log *log.Logger, prefix string) (c *Counter, err error) {
	return &Counter{
		log:    log,
		Prefix: prefix,
	}, nil
}

// Run - do processing
func (c *Counter) Run(ctx context.Context, input chan NginxEvent) error {
	c.log.Print("started counter thread")

	handledEventsMetric := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: c.Prefix,
		Subsystem: "nsm",
		Name:      "handled_events",
		Help:      "nginx handled events",
	}, []string{"host", "status", "category"})

	handleTimesMetric := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: c.Prefix,
		Subsystem: "nsm",
		Name:      "handle_times",
		Help:      "nginx handle time",
	}, []string{"host", "category"})

	if err := prometheus.Register(handledEventsMetric); err != nil {
		return err
	}

	if err := prometheus.Register(handleTimesMetric); err != nil {
		return err
	}

	for ctx.Err() == nil {
		for v := range input {
			handledEventsMetric.WithLabelValues(v.HttpHost, v.Status, v.Category).Inc()

			for _, times := range strings.Split(v.UpstreamResponseTime, " ") {
				tm, err := strconv.ParseFloat(times, 64)
				if err != nil {
					continue
				}

				handleTimesMetric.WithLabelValues(v.HttpHost, v.Category).Add(tm)
			}
		}
	}

	return ctx.Err()
}

package internal

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"log"
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

	handledEvents := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: c.Prefix,
		Subsystem: "nsm",
		Name:      "handled_events",
		Help:      "nginx handled events",
	}, []string{"host", "category"})

	if err := prometheus.Register(handledEvents); err != nil {
		return err
	}

	for ctx.Err() == nil {
		for v := range input {
			handledEvents.WithLabelValues(v.HttpHost, v.Category).Inc()
		}
	}

	return ctx.Err()
}

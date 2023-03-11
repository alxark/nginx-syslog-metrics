package internal

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	"gopkg.in/mcuadros/go-syslog.v2"
	_ "gopkg.in/mcuadros/go-syslog.v2"
	"log"
	"sync"
	"time"
)

var receiverMetric *prometheus.CounterVec
var receiverOnce sync.Once

type Receiver struct {
	SyslogPort      int
	SyslogQueueSize int
	Prefix          string
	log             *log.Logger
}

func NewReceiver(log *log.Logger, syslogPort int, syslogQueueSize int, prefix string) (*Receiver, error) {
	r := &Receiver{}

	r.log = log
	r.SyslogPort = syslogPort
	r.SyslogQueueSize = syslogQueueSize
	r.Prefix = prefix

	receiverOnce.Do(func() {
		receiverMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: r.Prefix,
			Subsystem: "nsm",
			Name:      "incoming",
			Help:      "nginx messages received",
		}, []string{"type"})

		prometheus.MustRegister(receiverMetric)
	})

	return r, nil
}

// Run - read messages from syslog
func (r *Receiver) Run(outputChannel chan SyslogMessage) (err error) {
	r.log.Print("starting syslog incoming service")

	channel := make(syslog.LogPartsChannel, r.SyslogQueueSize)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC3164)
	server.SetHandler(handler)

	if err = server.ListenUDP(fmt.Sprintf("%s:%d", "0.0.0.0", r.SyslogPort)); err != nil {
		return err
	}

	if err = server.Boot(); err != nil {
		return err
	}

	for m := range channel {
		receiverMetric.WithLabelValues("received").Inc()
		var content map[string]string

		if err = json.Unmarshal([]byte(m["content"].(string)), &content); err != nil {
			r.log.Printf("failed to decode syslog content: %s", err.Error())
			receiverMetric.WithLabelValues("broken").Inc()
			continue
		}

		sMsg := SyslogMessage{
			Time:    m["timestamp"].(time.Time),
			Host:    m["hostname"].(string),
			Message: content,
		}

		receiverMetric.WithLabelValues("sent").Inc()
		outputChannel <- sMsg
	}

	return nil
}

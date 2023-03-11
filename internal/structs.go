package internal

import "time"

type SyslogMessage struct {
	Time    time.Time
	Host    string
	Message map[string]string
}

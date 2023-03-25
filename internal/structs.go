package internal

import "time"

type SyslogMessage struct {
	Time    time.Time
	Host    string
	Message string
}

type NginxEvent struct {
	Mode                   string // access / error logs
	Category               string // nginx event category
	Frontend               string
	RequestId              string `json:"request_id"`
	RequestLength          string `json:"request_length"`
	RemoteAddr             string `json:"request_addr"`
	RemoteUser             string `json:"request_user"`
	RemotePort             string `json:"remote_port"`
	Request                string `json:"request"`
	RequestUri             string `json:"request_uri"`
	RequestMethod          string `json:"request_method"`
	Args                   string `json:"args"`
	Status                 string `json:"status"`
	BodyBytesSent          string `json:"body_bytes_sent"`
	BytesSent              string `json:"bytes_sent"`
	HttpHost               string `json:"http_host"`
	ServerName             string `json:"server_name"`
	Scheme                 string `json:"scheme"`
	RequestTime            string `json:"request_time"`
	UpstreamAddr           string `json:"upstream_addr"`
	UpstreamResponseTime   string `json:"upstream_response_time"`
	UpstreamResponseLength string `json:"upstream_response_length"`
}

type CategoriserConfig struct {
	SourceRegexp string `json:"source"`
	Target       string `json:"target"`
}

# nginx-syslog-metrics
Nginx Syslog Metrics - is a service designed to collect logs from Nginx HTTP servers using the Syslog protocol. 
The service efficiently transforms and exposes these logs as Prometheus metrics on a separate, user-specified port. 
Enjoy enhanced monitoring and analysis of your Nginx server performance with seamless integration into your existing Prometheus setup.

## How to setup

### Run nginx-syslog-metric

Use docker to run nginx-syslog-metric

```bash
docker build -t nginx-syslog-metric . 
docker run --name=nginx-syslog-metric -d -p 9090:9090 -p 8080:8080 -t -i nginx-syslog-metric
```

### Configure Nginx

Use following configuration to enable syslog logging in Nginx

```bash
http {
  
  log_format syslog escape=json '{'
    '"time_local":"$time_local",'
    '"remote_addr":"$remote_addr",'
    '"remote_user":"$remote_user",'
    '"status":"$status",'
    '"body_bytes_sent":"$body_bytes_sent",'
    '"http_referer":"$http_referer",'
    '"http_user_agent":"$http_user_agent",'
    '"http_x_forwarded_for":"$http_x_forwarded_for",'
    '"request_time":"$request_time",'
    '"upstream_response_time":"$upstream_response_time",'
    '"upstream_addr":"$upstream_addr",'
    '"upstream_status":"$upstream_status",'
    '"upstream_cache_status":"$upstream_cache_status",'
    '"upstream_bytes_sent":"$upstream_bytes_sent",'
    '"upstream_connect_time":"$upstream_connect_time",'
    '"upstream_header_time":"$upstream_header_time",'
    '"upstream_response_length":"$upstream_response_length",'
    '"upstream_response_time":"$upstream_response_time",'
    '"upstream_tries":"$upstream_tries",'
    '"scheme":"$scheme",'
    '"host":"$host",'
    '"request":"$request",'
    '"request_method":"$request_method",'
    '"request_uri":"$request_uri",'
    '"request_filename":"$request_filename",'
    '"request_uri":"$request_uri",'
    '"request_length":"$request_length",'
    '"request_time":"$request_time",'
    '"server_protocol":"$server_protocol",'
    '"server_port":"$server_port" '
    '}';
```

Add logging to required servers / locations

```nginx
   access_log syslog:server=localhost:9090 json_analytics;
   error_log  syslog:server=localhost:9090;
```

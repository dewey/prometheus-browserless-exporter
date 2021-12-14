# Prometheus Browserless Exporter

This is an exporter that exposes [Browserless](https://github.com/browserless/chrome) metrics for Prometheus. Browserless exposes [some](https://github.com/browserless/chrome/pull/381) Prometheus metrics already but not all of them.

The exporter provides access to the following additional metrics:

```
# HELP browserless_concurrent_max Maximum number of concurrent sessions observed
# TYPE browserless_concurrent_max gauge
browserless_concurrent_max 0
# HELP browserless_error Number of requests resulting in an error
# TYPE browserless_error counter
browserless_error 0
# HELP browserless_exporter_errors_total The total number of errors per collector
# TYPE browserless_exporter_errors_total counter
browserless_exporter_errors_total{collector="metrics-total"} 0
# HELP browserless_queued Number of requests that got queued
# TYPE browserless_queued counter
browserless_queued 0
# HELP browserless_rejected Number of requests that got rejected
# TYPE browserless_rejected counter
browserless_rejected 0
# HELP browserless_successful Number of successful requests
# TYPE browserless_successful counter
browserless_successful 0
# HELP browserless_time_max MaxTime as defined by browserless
# TYPE browserless_time_max gauge
browserless_time_max 0
# HELP browserless_time_mean MeanTime as defined by browserless
# TYPE browserless_time_mean gauge
browserless_time_mean 0
# HELP browserless_time_min MinTime as defined by browserless
# TYPE browserless_time_min gauge
browserless_time_min 0
# HELP browserless_time_total TotalTime as defined by browserless
# TYPE browserless_time_total gauge
browserless_time_total 0
# HELP browserless_timedout Number of timedout requests
# TYPE browserless_timedout counter
browserless_timedout 0
# HELP browserless_unhealthy Number of unhealthy requests
# TYPE browserless_unhealthy counter
browserless_unhealthy 0
```

# Usage

The following environment variables can be configured:

- `LISTEN_ADDR`: Address / Port for the exporter, default: *:3001*
- `TIMEOUT`: Timeout for the request to Browserless in seconds, default: *5*
- `DEBUG`: Debug mode and increased logging, default: false
- `BROWSERLESS_ENDPOINT`: The endpoint where Browserless exposes its metrics, default: *http://localhost:3000/metrics/total* 
- `METRICS_PATH`: Path where the Prometheus metrics are going to be exposed at, default: */metrics*

# Run

With Docker Compose:

```
docker-compose -f docker-compose.yml up -d
```

# Credits

Some inspiration was taken from these projects:

- JustWatch/[sql_exporter](https://github.com/justwatchcom/sql_exporter)
- metalmatze/[digitalocean_exporter](https://github.com/metalmatze/digitalocean_exporter)
# Prometheus Browserless Exporter

This is an exporter that exposes [Browserless](https://github.com/browserless/chrome) metrics for Prometheus. Browserless exposes [some](https://github.com/browserless/chrome/pull/381) Prometheus metrics already but not all of them.

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
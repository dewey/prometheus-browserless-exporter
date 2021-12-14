package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/peterbourgon/ff/v3"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	filterOption := level.AllowInfo()
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, filterOption)
	logger = log.With(logger,
		"ts", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)

	fs := flag.NewFlagSet("prometheus-browserless-exporter", flag.ExitOnError)
	var (
		listenAddr          = fs.String("listen-addr", "localhost:3001", "listen address")
		timeoutSeconds      = fs.Int("timeout", 5, "timeout in seconds")
		debug               = fs.Bool("debug", false, "log debug information")
		browserlessEndpoint = fs.String("browserless-endpoint", "http://localhost:3000/metrics/total", "browserless metrics endpoint")
		metricsPath         = fs.String("metrics-endpoint", "/metrics", "path where prometheus metrics are going to be exposed")
	)

	if err := ff.Parse(fs, os.Args[1:],
		ff.WithEnvVarNoPrefix(),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ff.PlainParser),
	); err != nil {
		level.Error(logger).Log("msg", "failed to parse flags", "err", err)
		os.Exit(1)
	}

	if *debug {
		filterOption = level.AllowDebug()
	}

	level.Info(logger).Log("msg", "starting prometheus-browserless-exporter")
	timeout := time.Duration(*timeoutSeconds) * time.Millisecond
	errors := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "browserless_exporter_errors_total",
		Help: "The total number of errors per collector",
	}, []string{"collector"})

	// Client used to call metrics endpoints
	httpClient := http.Client{
		Timeout: time.Duration(*timeoutSeconds) * time.Second,
	}

	r := prometheus.NewRegistry()
	r.MustRegister(errors)
	r.MustRegister(NewBrowserlessTotalCollector(logger, httpClient, errors, timeout, *browserlessEndpoint))

	http.Handle(*metricsPath,
		promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>Browserless Exporter</title></head>
			<body>
			<h1>Browserless Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(logger).Log("msg", "listening", "addr", *listenAddr)
	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		level.Error(logger).Log("msg", "http listenandserve error", "err", err)
		os.Exit(1)
	}
}

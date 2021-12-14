package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

// BrowserlessTotalCollector collects metrics about a browserless instance.
type BrowserlessTotalCollector struct {
	logger              log.Logger
	httpClient          http.Client
	errors              *prometheus.CounterVec
	timeout             time.Duration
	browserlessEndpoint string

	Successful    *prometheus.Desc
	Error         *prometheus.Desc
	Queued        *prometheus.Desc
	Rejected      *prometheus.Desc
	Unhealthy     *prometheus.Desc
	Timedout      *prometheus.Desc
	TotalTime     *prometheus.Desc
	MeanTime      *prometheus.Desc
	MaxTime       *prometheus.Desc
	MinTime       *prometheus.Desc
	MaxConcurrent *prometheus.Desc
}

// NewBrowserlessTotalCollector returns a new BrowserlessCollector.to get the total metrics of an instance
func NewBrowserlessTotalCollector(logger log.Logger, httpClient http.Client, errors *prometheus.CounterVec, timeout time.Duration, browserlessEndpoint string) *BrowserlessTotalCollector {
	errors.WithLabelValues("metrics-total").Add(0)

	return &BrowserlessTotalCollector{
		logger:              logger,
		errors:              errors,
		timeout:             timeout,
		browserlessEndpoint: browserlessEndpoint,
		httpClient:          httpClient,

		Successful: prometheus.NewDesc(
			"browserless_successful",
			"Number of successful requests",
			nil, nil,
		),
		Error: prometheus.NewDesc(
			"browserless_error",
			"Number of requests resulting in an error",
			nil, nil,
		),
		Queued: prometheus.NewDesc(
			"browserless_queued",
			"Number of requests that got queued",
			nil, nil,
		),
		Rejected: prometheus.NewDesc(
			"browserless_rejected",
			"Number of requests that got rejected",
			nil, nil,
		),
		Unhealthy: prometheus.NewDesc(
			"browserless_unhealthy",
			"Number of unhealthy requests",
			nil, nil,
		),
		Timedout: prometheus.NewDesc(
			"browserless_timedout",
			"Number of timedout requests",
			nil, nil,
		),
		TotalTime: prometheus.NewDesc(
			"browserless_time_total",
			"TotalTime as defined by browserless",
			nil, nil,
		),
		MeanTime: prometheus.NewDesc(
			"browserless_time_mean",
			"MeanTime as defined by browserless",
			nil, nil,
		),
		MaxTime: prometheus.NewDesc(
			"browserless_time_max",
			"MaxTime as defined by browserless",
			nil, nil,
		),
		MinTime: prometheus.NewDesc(
			"browserless_time_min",
			"MinTime as defined by browserless",
			nil, nil,
		),
		MaxConcurrent: prometheus.NewDesc(
			"browserless_concurrent_max",
			"Maximum number of concurrent sessions observed",
			nil, nil,
		),
	}
}

// Describe sends the super-set of all possible descriptors of metrics collected by this Collector.
func (c *BrowserlessTotalCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Successful
	ch <- c.Error
	ch <- c.Queued
	ch <- c.Rejected
	ch <- c.Unhealthy
	ch <- c.Timedout
	ch <- c.TotalTime
	ch <- c.MeanTime
	ch <- c.MaxTime
	ch <- c.MinTime
	ch <- c.MaxConcurrent
}

// Collect is called by the Prometheus registry when collecting metrics.
func (c *BrowserlessTotalCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := c.httpClient.Get(c.browserlessEndpoint)
	if err != nil {
		c.errors.WithLabelValues("metrics-total").Add(1)
		level.Warn(c.logger).Log("msg", "couldn't fetch metrics from total endpoint", "err", err)
		return
	}
	var m metricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		c.errors.WithLabelValues("metrics-total").Add(1)
		level.Warn(c.logger).Log("msg", "couldn't decode metrics response", "err", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		c.Successful,
		prometheus.CounterValue,
		float64(m.Successful),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Error,
		prometheus.CounterValue,
		float64(m.Error),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Queued,
		prometheus.CounterValue,
		float64(m.Queued),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Rejected,
		prometheus.CounterValue,
		float64(m.Rejected),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Unhealthy,
		prometheus.CounterValue,
		float64(m.Unhealthy),
	)
	ch <- prometheus.MustNewConstMetric(
		c.Timedout,
		prometheus.CounterValue,
		float64(m.Timedout),
	)
	ch <- prometheus.MustNewConstMetric(
		c.TotalTime,
		prometheus.GaugeValue,
		m.TotalTime,
	)
	ch <- prometheus.MustNewConstMetric(
		c.MeanTime,
		prometheus.GaugeValue,
		m.MeanTime,
	)
	ch <- prometheus.MustNewConstMetric(
		c.MaxTime,
		prometheus.GaugeValue,
		m.MaxTime,
	)
	ch <- prometheus.MustNewConstMetric(
		c.MinTime,
		prometheus.GaugeValue,
		m.MinTime,
	)
	ch <- prometheus.MustNewConstMetric(
		c.MaxConcurrent,
		prometheus.GaugeValue,
		float64(m.MaxConcurrent),
	)
}

// metricsResponse is the response returned by Browserless for an instance
type metricsResponse struct {
	Successful    int     `json:"successful"`
	Error         int     `json:"error"`
	Queued        int     `json:"queued"`
	Rejected      int     `json:"rejected"`
	Unhealthy     int     `json:"unhealthy"`
	Timedout      int     `json:"timedout"`
	TotalTime     float64 `json:"totalTime"`
	MeanTime      float64 `json:"meanTime"`
	MaxTime       float64 `json:"maxTime"`
	MinTime       float64 `json:"minTime"`
	MaxConcurrent int     `json:"maxConcurrent"`
}

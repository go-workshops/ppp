// Package metrics provides a simple wrapper around Prometheus metrics
// with easy to create and reuse metrics helper functions.
// Check out the examples file, for a more detailed list of various Prometheus metrics and how to use them:
// https://github.com/prometheus/client_golang/blob/main/prometheus/examples_test.go
package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const appNameLabel = "app_name"

var (
	// DefaultPrefix is the default prefix used for all metric names.
	// i.e: DefaultPrefix_your_metric_name
	DefaultPrefix = "ppp"

	// DefaultProvider represents the default metrics provider.
	DefaultProvider Provider = NewPrometheusProvider(PrometheusProviderOpts{})

	// ConstLabels are the default const labels applied all newly registered metrics.
	// Make sure to set all the necessary labels (on application setup) before making use of or creating any metrics.
	ConstLabels = ConstMetricLabels{labels: map[string]string{}}
)

var (
	mu         sync.Mutex
	counters   = map[string]CounterVecMetric{}
	gauges     = map[string]GaugeVecMetric{}
	histograms = map[string]ObserverVecMetric{}
	summaries  = map[string]ObserverVecMetric{}
)

// SetAppName sets the application name for the default metrics provider.
// This will create a const label with the key AppNameLabel for every registered metric.
func SetAppName(appName string) {
	ConstLabels.Set(appNameLabel, appName)
}

// GetAppName returns the application name for the default metrics provider.
func GetAppName() string {
	appName, _ := ConstLabels.Get(appNameLabel)
	return appName
}

// CounterVecMetric represents a vector counter metric containing a variation
// of the same metric under different labels.
type CounterVecMetric interface {
	With(labels map[string]string) CounterMetric
}

// CounterMetric represents a counter metric.
type CounterMetric interface {
	Inc()
	Add(float64)
}

// GaugeVecMetric represents a vector gauge metric containing a variation
// of the same metric under different labels.
type GaugeVecMetric interface {
	With(labels map[string]string) GaugeMetric
}

// GaugeMetric represents a gauge metric.
type GaugeMetric interface {
	Set(float64)
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
	SetToCurrentTime()
}

// ObserverVecMetric represents a vector observer(histogram/summary) metric containing a variation
// of the same metric under different labels.
type ObserverVecMetric interface {
	With(labels map[string]string) ObserverMetric
}

// ObserverMetric represents a Histogram / Summary metric.
type ObserverMetric interface {
	Observe(float64)
}

// Provider represents a metric provider, i.e: Prometheus.
type Provider interface {
	NewCounter(name, help string, constLabels map[string]string, labels ...string) CounterVecMetric
	NewGauge(name, help string, constLabels map[string]string, labels ...string) GaugeVecMetric
	NewHistogram(name, help string, constLabels map[string]string, buckets []float64, labels ...string) ObserverVecMetric
	NewSummary(name, help string, constLabels map[string]string, objectives map[float64]float64, labels ...string) ObserverVecMetric
	WithCollector(collector prometheus.Collector) Provider
}

// Counter creates or references an existing counter metric.
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a CounterMetric.
func Counter(name string, args ...string) CounterMetric {
	return counter(name, help(args), ConstLabels.labels).With(map[string]string{})
}

// CounterVec creates or references an existing counter vector metric.
// Use this function instead, if you plan on dynamically adding custom labels
// to the CounterMetric, which involves an extra step of calling
// .With(map[string]string{"label_name": "label_value"}), which then
// gives the caller access to a CounterMetric to work with.
func CounterVec(name string, args ...string) CounterVecMetric {
	return counter(name, help(args), ConstLabels.labels, labels(args)...)
}

func counter(name string, help string, constLabels map[string]string, labels ...string) CounterVecMetric {
	mu.Lock()
	defer mu.Unlock()

	name = fqdn(name)
	c, ok := counters[name]
	if ok {
		return c
	}

	c = DefaultProvider.NewCounter(name, help, constLabels, labels...)
	counters[name] = c
	return c
}

// Gauge creates or references an existing gauge metric.
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a GaugeMetric.
func Gauge(name string, args ...string) GaugeMetric {
	return gauge(name, help(args), ConstLabels.labels).With(map[string]string{})
}

// GaugeVec creates or references an existing gauge vector metric.
// Use this function instead, if you plan on dynamically adding custom labels
// to the GaugeMetric, which involves an extra step of calling
// .With(map[string]string{"label_name": "label_value"}), which then
// gives the caller access to a GaugeMetric to work with.
func GaugeVec(name string, args ...string) GaugeVecMetric {
	return gauge(name, help(args), ConstLabels.labels, labels(args)...)
}

func gauge(name string, help string, constLabels map[string]string, labels ...string) GaugeVecMetric {
	mu.Lock()
	defer mu.Unlock()

	name = fqdn(name)
	g, ok := gauges[name]
	if ok {
		return g
	}

	g = DefaultProvider.NewGauge(name, help, constLabels, labels...)
	gauges[name] = g
	return g
}

// Histogram creates or references an existing histogram metric.
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a ObserverMetric (histogram).
func Histogram(name string, args ...string) ObserverMetric {
	return histogram(name, help(args), ConstLabels.labels, []float64{}).With(map[string]string{})
}

// HistogramWithBuckets creates or references an existing histogram metric with custom buckets.
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a ObserverMetric (histogram), and is initialized with custom buckets.
func HistogramWithBuckets(name string, buckets []float64, args ...string) ObserverMetric {
	return histogram(name, help(args), ConstLabels.labels, buckets).With(map[string]string{})
}

// HistogramVec creates or references an existing histogram vector metric.
// Use this function instead, if you plan on dynamically adding custom labels
// to the ObserverMetric (histogram), which involves an extra step of calling
// .With(map[string]string{"label_name": "label_value"}), which then
// gives the caller access to a ObserverMetric (histogram) to work with.
func HistogramVec(name string, args ...string) ObserverVecMetric {
	return histogram(name, help(args), ConstLabels.labels, []float64{}, labels(args)...)
}

// HistogramVecWithBuckets creates or references an existing histogram vector metric with custom buckets.
// Use this function instead, if you plan on dynamically adding custom labels
// to the ObserverMetric (histogram), which involves an extra step of calling
// .With(map[string]string{"label_name": "label_value"}), which then
// gives the caller access to a ObserverMetric (histogram) to work with and is initialized with custom buckets..
func HistogramVecWithBuckets(name string, buckets []float64, args ...string) ObserverVecMetric {
	return histogram(name, help(args), ConstLabels.labels, buckets, labels(args)...)
}

func histogram(name string, help string, constLabels map[string]string, buckets []float64, labels ...string) ObserverVecMetric {
	mu.Lock()
	defer mu.Unlock()

	name = fqdn(name)
	h, ok := histograms[name]
	if ok {
		return h
	}

	h = DefaultProvider.NewHistogram(name, help, constLabels, buckets, labels...)
	histograms[name] = h
	return h
}

// Summary creates or references an existing summary metric.
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a ObserverMetric (summary).
func Summary(name string, args ...string) ObserverMetric {
	return summary(name, help(args), ConstLabels.labels, map[float64]float64{}).With(map[string]string{})
}

// SummaryWithObjectives creates or references an existing summary metric with objectives → map[quantile:absolute error].
// Use this function, if the metric does not have any custom dynamic labels,
// which also gives the caller direct access to a ObserverMetric (summary),
// which is also initialized with objectives → map[quantile:absolute error].
// For more information about quantiles (objectives) check out:
// https://en.wikipedia.org/wiki/Quantile
// https://en.wikipedia.org/wiki/Percentile
func SummaryWithObjectives(name string, objectives map[float64]float64, args ...string) ObserverMetric {
	return summary(name, help(args), ConstLabels.labels, objectives).With(map[string]string{})
}

// SummaryVec creates or references an existing summary vector metric.
// Use this function instead, if you plan on dynamically adding custom labels
// to the ObserverMetric (summary), which involves an extra step of calling
// .With(map[string]string{"label_name": "label_value"}), which then
// gives the caller access to a ObserverMetric (summary) to work with.
func SummaryVec(name string, args ...string) ObserverVecMetric {
	return summary(name, help(args), ConstLabels.labels, map[float64]float64{}, labels(args)...)
}

// SummaryVecWithObjectives creates or references an existing summary vector metric
// with objectives → map[quantile:absolute error]. Use this function, if the metric does not have
// any custom dynamic labels, which also gives the caller direct access to a ObserverMetric (summary),
// which is also initialized with objectives → map[quantile:absolute error].
// For more information about quantiles (objectives) check out:
// https://en.wikipedia.org/wiki/Quantile
// https://en.wikipedia.org/wiki/Percentile
func SummaryVecWithObjectives(name string, objectives map[float64]float64, args ...string) ObserverVecMetric {
	return summary(name, help(args), ConstLabels.labels, objectives, labels(args)...)
}

func summary(name string, help string, constLabels map[string]string, objectives map[float64]float64, labels ...string) ObserverVecMetric {
	mu.Lock()
	defer mu.Unlock()

	name = fqdn(name)
	s, ok := summaries[name]
	if ok {
		return s
	}

	s = DefaultProvider.NewSummary(name, help, constLabels, objectives, labels...)
	summaries[name] = s
	return s
}

// help extracts the metric help message from a variadic list of fields
func help(args []string) string {
	h := ""
	if len(args) > 0 {
		h = args[0]
	}
	return h
}

// labels extracts the metric labels from a variadic list of fields
func labels(args []string) []string {
	if len(args) < 2 {
		return []string{}
	}
	return args[1:]
}

func fqdn(name string) string {
	return DefaultPrefix + "_" + name
}

// RegisterCollector registers a collector with the default provider.
func RegisterCollector(collector prometheus.Collector) {
	DefaultProvider.WithCollector(collector)
}

// ConstMetricLabels represents the constant metric labels wrapper.
type ConstMetricLabels struct {
	labels map[string]string
	mu     sync.RWMutex
}

// Set sets a constant metric label that will be available to all registered metrics.
// If the value is empty, the const label will be omitted.
func (c *ConstMetricLabels) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value != "" {
		c.labels[key] = value
	}
}

// Get gets a constant metric label available on all registered metrics.
func (c *ConstMetricLabels) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.labels[key]

	return val, ok
}

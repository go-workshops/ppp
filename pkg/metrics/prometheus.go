package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry                               = prometheus.NewRegistry()
	registerer       prometheus.Registerer = registry
	gatherer         prometheus.Gatherer   = registry
	metricCollectors                       = newCollectors
)

func init() {
	registerer = prometheus.NewRegistry()
	registerer = prometheus.WrapRegistererWithPrefix(
		fqdn(""),
		prometheus.WrapRegistererWith(ConstLabels.labels, registry),
	)
	registerer.MustRegister(metricCollectors()...)
}

// PrometheusProviderOpts represents the Prometheus metrics configuration options.
type PrometheusProviderOpts struct {
	prometheus.Registerer
	prometheus.Gatherer
}

// NewPrometheusProvider creates a new Prometheus provider that implements Provider using Prometheus metrics.
func NewPrometheusProvider(opts PrometheusProviderOpts) PrometheusProvider {
	if opts.Registerer != nil {
		registerer = opts.Registerer
	}
	if opts.Gatherer != nil {
		gatherer = opts.Gatherer
	}
	p := PrometheusProvider{
		registerer: registerer,
		gatherer:   gatherer,
	}
	return p
}

// PrometheusProvider represents the implementation for Prometheus provider.
type PrometheusProvider struct {
	registerer prometheus.Registerer
	gatherer   prometheus.Gatherer
}

// NewCounter creates a new Prometheus counter vector metric.
func (p PrometheusProvider) NewCounter(name, help string, constLabels map[string]string, labels ...string) CounterVecMetric {
	vec := promauto.With(p.registerer).NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Help:        help,
			ConstLabels: constLabels,
		},
		labels,
	)
	return counterVec{vec}
}

// counterVec represents an internal counter vec type that implements CounterVecMetric
type counterVec struct {
	*prometheus.CounterVec
}

func (c counterVec) With(labels map[string]string) CounterMetric {
	return c.CounterVec.With(labels)
}

// NewGauge creates a new Prometheus gauge vector metric.
func (p PrometheusProvider) NewGauge(name, help string, constLabels map[string]string, labels ...string) GaugeVecMetric {
	vec := promauto.With(p.registerer).NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        help,
			ConstLabels: constLabels,
		},
		labels,
	)
	return gaugeVec{vec}
}

// gaugeVec represents an internal gauge vec type that implements GaugeVecMetric
type gaugeVec struct {
	*prometheus.GaugeVec
}

func (g gaugeVec) With(labels map[string]string) GaugeMetric {
	return g.GaugeVec.With(labels)
}

// NewHistogram creates a new Prometheus histogram vector metric.
func (p PrometheusProvider) NewHistogram(name, help string, constLabels map[string]string, buckets []float64, labels ...string) ObserverVecMetric {
	vec := promauto.With(p.registerer).NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        help,
			ConstLabels: constLabels,
			Buckets:     buckets,
		},
		labels,
	)
	return histogramVec{vec}
}

// histogramVec represents an internal histogram vec type that implements ObserverVecMetric
type histogramVec struct {
	*prometheus.HistogramVec
}

func (h histogramVec) With(labels map[string]string) ObserverMetric {
	return h.HistogramVec.With(labels)
}

// NewSummary creates a new Prometheus summary vector metric.
func (p PrometheusProvider) NewSummary(name, help string, constLabels map[string]string, objectives map[float64]float64, labels ...string) ObserverVecMetric {
	vec := promauto.With(p.registerer).NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        name,
			Help:        help,
			ConstLabels: constLabels,
			Objectives:  objectives,
		},
		labels,
	)
	return summaryVec{vec}
}

// summaryVec represents an internal summary vec type that implements ObserverVecMetric
type summaryVec struct {
	*prometheus.SummaryVec
}

func (s summaryVec) With(labels map[string]string) ObserverMetric {
	return s.SummaryVec.With(labels)
}

// PrometheusHandler creates a new http.Handler that exposes Prometheus metrics over HTTP.
func PrometheusHandler() http.Handler {
	return promhttp.InstrumentMetricHandler(
		registerer,
		promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}),
	)
}

// WithCollector registers a new collector with the Prometheus provider.
func (p PrometheusProvider) WithCollector(collector prometheus.Collector) Provider {
	p.registerer.Unregister(collector)
	p.registerer.MustRegister(collector)
	return p
}

func newCollectors() []prometheus.Collector {
	return []prometheus.Collector{
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	}
}

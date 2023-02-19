package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/src/main/app/config"
	"log"
)

type IMetricCollector interface {
	IncrementCounter(name string, tags ...string)
	RecordExecutionTime(name string, value int64, tags ...string)
}

var Collector = newMetricsCollector()
var counters = make(map[string]prometheus.Counter)
var summaries = make(map[string]prometheus.Summary)
var genericCounter *prometheus.CounterVec

type metricsCollector struct {
}

func newMetricsCollector() *metricsCollector {

	success := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "consumers",
			Name:        "consumers_pusher_success",
			Help:        "How many messages processed.",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
	)
	prometheus.MustRegister(success)
	counters["consumers.pusher.success"] = success

	errors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "consumers",
			Name:        "consumers_pusher_error",
			Help:        "How many messages can't be processed.",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
	)
	prometheus.MustRegister(errors)
	counters["consumers.pusher.errors"] = errors

	pusher20x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "consumers",
			Name:        "consumers_pusher_http_200",
			Help:        "How many ACK.",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
	)
	prometheus.MustRegister(pusher20x)
	counters["consumers.pusher.http.20x"] = pusher20x

	pusher40x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "consumers",
			Name:        "consumers_pusher_http_40x",
			Help:        "How many messages can't be processed.",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
	)
	prometheus.MustRegister(pusher40x)
	counters["consumers.pusher.http.40x"] = pusher40x

	pusher50x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   "consumers",
			Name:        "consumers_pusher_http_50x",
			Help:        "How many messages can't be processed.",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
	)
	prometheus.MustRegister(pusher50x)
	counters["consumers.pusher.http.50x"] = pusher50x

	client := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   "consumers",
		Name:        "consumers_pusher_http_time",
		Help:        "Duration of the login request.",
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
	})
	prometheus.MustRegister(client)
	summaries["consumers.pusher.http.time"] = client

	generic := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "consumers_pusher_generic_counter",
			Help:        "Generic counter. Filtered by name",
			ConstLabels: prometheus.Labels{"env": config.String("app.env"), "app": config.String("app.name")},
		},
		[]string{"name"},
	)
	prometheus.MustRegister(generic)
	genericCounter = generic

	return &metricsCollector{}
}

func (m metricsCollector) IncrementCounter(name string) {
	if counter, ok := counters[name]; ok {
		counter.Inc()
	} else {
		genericCounter.WithLabelValues(name).Inc()
	}
}

func (m metricsCollector) RecordExecutionTime(name string, value int64) {
	if summary, ok := summaries[name]; ok {
		summary.Observe(float64(value))
	} else {
		log.Printf("missing metric collector: %s", name)
	}
}

package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/src/main/app/config"
	"github.com/src/main/app/log"
	"github.com/ugurcsen/gods-generic/maps/hashmap"
)

type IMetricCollector interface {
	IncrementCounter(name Name)
	Record(name Name, value int)
	RecordExecutionTime(name Name, value time.Duration)
}

type Name string

// Pusher metrics.
const (
	PusherSuccess     Name = "app_pusher_success"
	PusherError       Name = "app_pusher_error"
	PusherStatusOK    Name = "app_pusher_http_200"
	PusherStatus40x   Name = "app_pusher_http_4xx"
	PusherStatus50x   Name = "app_pusher_http_5xx"
	PusherHTTPTime    Name = "app_pusher_http_time"
	PusherHTTPTimeout Name = "app_pusher_http_timeout"
	Generic           Name = "app_pusher_generic_counter"
)

// Consumer and queue metrics.
const (
	ApproximateNumberOfMessages Name = "app_approximate_number_of_messages"
	CurrentWorkers              Name = "app_current_workers"
)

var (
	Collector         = newMetricsCollector()
	counters          = hashmap.New[Name, prometheus.Counter]()
	summaries         = hashmap.New[Name, prometheus.Summary]()
	genericCounter    *prometheus.CounterVec
	namespace, labels = "consumers", prometheus.Labels{
		"env":   config.String("app.env"),
		"app":   config.String("app.name"),
		"scope": config.String("SCOPE"),
	}
)

type metricsCollector struct {
}

func newMetricsCollector() *metricsCollector {
	pusherSuccess := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherSuccess),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusherSuccess)
	counters.Put(PusherSuccess, pusherSuccess)

	pusherErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherError),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusherErrors)
	counters.Put(PusherError, pusherErrors)

	pusher20x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatusOK),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher20x)
	counters.Put(PusherStatusOK, pusher20x)

	pusher40x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatus40x),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher40x)
	counters.Put(PusherStatus40x, pusher40x)

	pusher50x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatus50x),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher50x)
	counters.Put(PusherStatus50x, pusher50x)

	approximateNumberOfMessages := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   namespace,
			Name:        string(ApproximateNumberOfMessages),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(approximateNumberOfMessages)
	summaries.Put(ApproximateNumberOfMessages, approximateNumberOfMessages)

	currentWorkers := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Namespace:   namespace,
			Name:        string(CurrentWorkers),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(currentWorkers)
	summaries.Put(CurrentWorkers, currentWorkers)

	client := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   namespace,
		Name:        string(PusherHTTPTime),
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		ConstLabels: labels,
	})
	prometheus.MustRegister(client)
	summaries.Put(PusherHTTPTime, client)

	pusherTimeout := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherHTTPTimeout),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusherTimeout)
	counters.Put(PusherHTTPTimeout, pusherTimeout)

	generic := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        string(Generic),
			Help:        "Generic counter. This is a fallback. Review your and metric.",
			ConstLabels: labels,
		},
		[]string{"name"},
	)
	prometheus.MustRegister(generic)
	genericCounter = generic

	return &metricsCollector{}
}

func (m metricsCollector) IncrementCounter(name Name) {
	if counter, ok := counters.Get(name); ok {
		counter.Inc()
	} else {
		log.Warnf("missing metric collector %s, fallback to generic metric collector", name)
		genericCounter.WithLabelValues(string(name)).Inc()
	}
}

func (m metricsCollector) Record(name Name, value int) {
	if summary, ok := summaries.Get(name); ok {
		summary.Observe(float64(value))
	} else {
		log.Warnf("missing metric collector: %s", string(name))
	}
}

func (m metricsCollector) RecordExecutionTime(name Name, value time.Duration) {
	if summary, ok := summaries.Get(name); ok {
		elapsedTime := float64(value.Nanoseconds()) / 1e9
		summary.Observe(elapsedTime)
	} else {
		log.Warnf("missing time metric collector: %s", string(name))
	}
}

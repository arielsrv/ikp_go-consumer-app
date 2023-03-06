package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/src/main/app/config"
	"github.com/src/main/app/log"
)

type IMetricCollector interface {
	IncrementCounter(name Name)
	RecordExecutionTime(name Name, value int64)
}

type Name string

const (
	PusherSuccess     Name = "pusher_success"
	PusherError       Name = "pusher_error"
	PusherStatusOK    Name = "pusher_http_200"
	PusherStatus40x   Name = "pusher_http_4xx"
	PusherStatus50x   Name = "pusher_http_5xx"
	PusherHTTPTime    Name = "pusher_http_time"
	PusherHTTPTimeout Name = "pusher_http_timeout"
	Generic           Name = "pusher_generic_counter"
)

var (
	Collector         = newMetricsCollector()
	counters          = make(map[string]prometheus.Counter)
	summaries         = make(map[string]prometheus.Summary)
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
	counters[string(PusherSuccess)] = pusherSuccess

	pusherErrors := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherError),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusherErrors)
	counters[string(PusherError)] = pusherErrors

	pusher20x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatusOK),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher20x)
	counters[string(PusherStatusOK)] = pusher20x

	pusher40x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatus40x),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher40x)
	counters[string(PusherStatus40x)] = pusher40x

	pusher50x := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherStatus50x),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusher50x)
	counters[string(PusherStatus50x)] = pusher50x

	client := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   namespace,
		Name:        string(PusherHTTPTime),
		Objectives:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		ConstLabels: labels,
	})
	prometheus.MustRegister(client)
	summaries[string(PusherHTTPTime)] = client

	pusherTimeout := prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace:   namespace,
			Name:        string(PusherHTTPTimeout),
			ConstLabels: labels,
		},
	)
	prometheus.MustRegister(pusherTimeout)
	counters[string(PusherHTTPTimeout)] = pusherTimeout

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
	if counter, ok := counters[string(name)]; ok {
		counter.Inc()
	} else {
		log.Warnf("missing metric collector %s, fallback to generic metric collector", name)
		genericCounter.WithLabelValues(string(name)).Inc()
	}
}

func (m metricsCollector) RecordExecutionTime(name Name, value time.Duration) {
	if summary, ok := summaries[string(name)]; ok {
		elapsedTime := float64(value.Nanoseconds()) / 1e9
		summary.Observe(elapsedTime)
	} else {
		log.Warnf("missing time metric collector: %s", string(name))
	}
}

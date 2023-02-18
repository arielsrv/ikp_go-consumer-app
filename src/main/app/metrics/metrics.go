package metrics

import "log"

type IMetricCollector interface {
	IncrementCounter(tag string)
}

type MetricCollector struct {
}

func (m MetricCollector) IncrementCounter(tag string) {
	log.Println("metric tag: increment counter")
}

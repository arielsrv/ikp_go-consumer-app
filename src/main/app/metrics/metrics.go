package metrics

type IMetricCollector interface {
	IncrementCounter(tag string)
}

var Collector = &metricsCollector{}

type metricsCollector struct {
}

func (m metricsCollector) IncrementCounter(name string, tags ...string) {
	// @TODO invoke real metric collector
}

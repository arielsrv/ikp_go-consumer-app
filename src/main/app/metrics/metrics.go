package metrics

type IMetricCollector interface {
	IncrementCounter(tag string)
}

type MetricCollector struct {
}

func (m MetricCollector) IncrementCounter(tag string) {
	// @TODO: increment metric for delivery notifications
}

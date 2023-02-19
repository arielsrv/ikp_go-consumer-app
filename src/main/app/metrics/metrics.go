package metrics

type IMetricCollector interface {
	IncrementCounter(name string, tags ...string)
	RecordExecutionTime(name string, value int64, tags ...string)
}

var Collector = &metricsCollector{}

type metricsCollector struct {
}

func (m metricsCollector) IncrementCounter(string, ...string) {
	//TODO implement me
}

func (m metricsCollector) RecordExecutionTime(string, int64, ...string) {
	//TODO implement me
}

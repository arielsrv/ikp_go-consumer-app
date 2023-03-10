package metrics_test

import (
	"testing"

	"github.com/src/main/app/metrics"
)

func TestMetricsCollector_IncrementCounter(t *testing.T) {
	metrics.Collector.IncrementCounter(metrics.PusherSuccess)
	metrics.Collector.IncrementCounter("fallback")
}
func TestMetricsCollector_RecordExecutionTime(t *testing.T) {
	metrics.Collector.RecordExecutionTime(metrics.PusherHTTPTime, 2000)
	metrics.Collector.RecordExecutionTime("fallback", 2000)
}

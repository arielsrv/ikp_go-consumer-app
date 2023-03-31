package metrics_test

import (
	"testing"

	"github.com/src/main/app/metrics"
)

func TestMetricsCollector_IncrementCounter(t *testing.T) {
	metrics.Collector.IncrementCounter(metrics.PusherSuccess)
	metrics.Collector.IncrementCounter("fallback")
	t.Log("done")
}

func TestMetricsCollector_RecordE(t *testing.T) {
	metrics.Collector.Record(metrics.CurrentWorkers, 2000)
	metrics.Collector.Record("fallback", 2000)
	t.Log("done")
}

func TestMetricsCollector_RecordExecutionTime(t *testing.T) {
	metrics.Collector.RecordExecutionTime(metrics.PusherHTTPTime, 2000)
	metrics.Collector.RecordExecutionTime("fallback", 2000)
	t.Log("done")
}

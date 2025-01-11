package testutil

import (
	prom_pb "github.com/prometheus/client_model/go"
	//"github.com/davecgh/go-spew/spew"
)

type LabelTests map[string]string
type MetricTests map[string]struct {
	Labels LabelTests
	Value  float64
}

func CompareMetrics(metrics []*prom_pb.MetricFamily, compare MetricTests) bool {
	num_ok := 0
	for i := 0; i < len(metrics); i++ {
		// Check if metric exists in compare sample
		key := metrics[i].GetName()
		test, ok := compare[key]
		if !ok {
			continue
		}
		metric := metrics[i].GetMetric()[0]

		// Compare labels against .Labels in the sample
		labels_ok := 0
		labels := metric.GetLabel()
		for i := 0; i < len(labels); i++ {
			label := labels[i]
			r, ok := test.Labels[label.GetName()]
			if !ok { // Skip label if it does not exist in sample
				continue
			}

			if label.GetValue() == r {
				labels_ok++
			}
		}
		if labels_ok != len(labels) {
			continue
		}

		// Compare value against .Value
		var val_m float64
		switch metrics[i].GetType() {
		case prom_pb.MetricType_COUNTER:
			val_m = metric.GetCounter().GetValue()
		case prom_pb.MetricType_GAUGE:
			val_m = metric.GetGauge().GetValue()
		default:
			continue
		}

		if test.Value == val_m {
			num_ok++
		}
	}
	return num_ok == len(compare)
}

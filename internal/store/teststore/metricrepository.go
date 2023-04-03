package teststore

import (
	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type Counter int64
type Gauge float64
type GaugeMap map[string]Gauge
type CounterMap map[string]Counter

type MetricRepository struct {
	gaugeMap   GaugeMap
	counterMap CounterMap
	metrics    map[string]*model.Metric
}

// UpdateMetric ...
func (mr *MetricRepository) UpdateMetric(m *model.Metric) error {

	// if m.MType == "gauge" {
	// 	mr.metrics[m.Name] = m.Value
	// } else if m.MType == "counter" {
	// 	// mr.metrics[m.Name] += *m.Delta // = mr.metrics[m.Name] + m.Delta
	// } else {
	// 	return errors.New("unsupported metric type")
	// }
	// mr.metrics[m.Name] = m
	return nil
}

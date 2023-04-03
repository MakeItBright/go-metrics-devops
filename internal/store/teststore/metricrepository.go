package teststore

type GaugeMap map[string]float64 // новое значение должно замещать предыдущее.
type CounterMap map[string]int64 // новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.

type MetricRepository struct {
	gaugeMap   GaugeMap
	counterMap CounterMap
}

func (mr *MetricRepository) SaveCounterValue(name string, counter int64) {

}
func (mr *MetricRepository) SaveGaugeValue(name string, gauge float64) {

}
func (mr *MetricRepository) GetCounterValue(name string) (int64, error) {
	return 0, nil
}
func (mr *MetricRepository) GetGaugeValue(name string) (float64, error) {
	return 0, nil
}
func (mr *MetricRepository) GetAllValues() string {

	return ""
}

// UpdateMetric ...
// func (mr *MetricRepository) UpdateMetric(m *model.Metric) error {

// 	// if m.MType == "gauge" {
// 	// 	mr.metrics[m.Name] = m.Value
// 	// } else if m.MType == "counter" {
// 	// 	// mr.metrics[m.Name] += *m.Delta // = mr.metrics[m.Name] + m.Delta
// 	// } else {
// 	// 	return errors.New("unsupported metric type")
// 	// }
// 	// mr.metrics[m.Name] = m
// 	return nil
// }

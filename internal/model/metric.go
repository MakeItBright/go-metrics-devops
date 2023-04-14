package model

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type MetricName string

type MetricPath string
type Metric struct {
	Name  MetricName // имя метрики
	Type  MetricType // параметр, принимающий значение gauge или counter
	Delta int64      // значение метрики в случае передачи counter
	Value float64    // значение метрики в случае передачи gauge
}

func (m *Metric) Validate() error {
	return nil
}

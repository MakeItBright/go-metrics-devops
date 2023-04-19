package model

import (
	"fmt"
)

type MetricName string // тип для имени метрики

type MetricType string // тип для типа метрики

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type Metric struct {
	Name  MetricName // имя метрики
	Type  MetricType // параметр, принимающий значение gauge или counter
	Delta int64      // значение метрики в случае передачи counter
	Value float64    // значение метрики в случае передачи gauge
}

func (m *Metric) Validate() error {
	return nil
}

// GetValue возвращает значение метрики в виде строки
func (m *Metric) GetValue() string {
	switch m.Type {
	case MetricTypeGauge:
		return fmt.Sprintf("%v", m.Value)
	case MetricTypeCounter:
		return fmt.Sprintf("%v", m.Delta)
	default:
		return ""
	}
}

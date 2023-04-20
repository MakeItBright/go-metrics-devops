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
	Name  MetricName `json:"id"`              // имя метрики
	Type  MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64      `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64    `json:"value,omitempty"` // значение метрики в случае передачи gauge
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

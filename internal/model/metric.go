package model

const (
	Counter = "counter" // новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
	Gauge   = "gauge"   // новое значение должно замещать предыдущее.
)

type Metric struct {
	Name  string  // имя метрики
	MType string  // параметр, принимающий значение gauge или counter
	Delta int64   // значение метрики в случае передачи counter
	Value float64 // значение метрики в случае передачи gauge
}

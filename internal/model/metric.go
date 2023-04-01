package model

const (
	Counter = "counter" // новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
	Gauge   = "gauge"   // новое значение должно замещать предыдущее.
)

type Metric struct {
	Name string
	Type string
}

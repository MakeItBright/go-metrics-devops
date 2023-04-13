package dto

type MetricRequest struct {
	Name  string  `` // имя метрики
	Type  string  // параметр, принимающий значение gauge или counter
	Delta int64   // значение метрики в случае передачи counter
	Value float64 // значение метрики в случае передачи gauge
}

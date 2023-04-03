package store

type MetricRepository interface {
	SaveCounterValue(name string, counter int64)
	SaveGaugeValue(name string, gauge float64)
	GetCounterValue(name string) (int64, error)
	GetGaugeValue(name string) (float64, error)
	GetAllValues() string
}

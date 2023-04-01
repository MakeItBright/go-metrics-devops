package store

type Store interface {
	Metric() MetricRepository
}

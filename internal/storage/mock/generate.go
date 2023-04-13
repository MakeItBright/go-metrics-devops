package mock

//go:generate mockgen -source ../metric.go -destination metric.go -package mock -mock_names Metric=MockMetricStorage

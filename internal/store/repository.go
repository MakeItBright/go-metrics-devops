package store

import (
	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type MetricRepository interface {
	UpdateMetric(*model.Metric) error
	// SaveGauge(name string, value int64) error
	// SaveCounter ()
}

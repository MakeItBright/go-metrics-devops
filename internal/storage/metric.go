package storage

import (
	"context"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type Metric interface {
	MetricStore(context.Context, model.Metric) error
	MetricFetch(context.Context, model.MetricType, model.MetricName) (model.Metric, error)
}

// type Storage interface {
// 	MetricInc(Metric, counter int64)
// 	MetricUpdate(string, gauge float64)
// 	GetCounterValue(name string) (int64, error)
// 	GetGaugeValue(name string) (float64, error)
// 	GetAllValues() string
// }

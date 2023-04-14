package storage

import (
	"context"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type Metric interface {
	MetricStore(context.Context, model.Metric) error
	MetricFetch(context.Context, model.MetricType, model.MetricName) (model.Metric, error)
	MetricAll() map[model.MetricPath]model.Metric
}

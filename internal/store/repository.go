package store

import (
	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type MetricRepository interface {
	Save(*model.Metric) error
}

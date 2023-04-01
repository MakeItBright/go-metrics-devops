package teststore

import (
	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type MetricRepository struct {
	metrics map[string]*model.Metric
}

// Save ...
func (mr *MetricRepository) Save(m *model.Metric) error {

	mr.metrics[m.Name] = m

	return nil
}

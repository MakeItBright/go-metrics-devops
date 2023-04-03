package teststore

import (
	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/store"
)

// Store ...
type Store struct {
	metricRepository *MetricRepository
}

func New() *Store {
	return &Store{}
}

// Metric ...
func (s *Store) Metric() store.MetricRepository {

	if s.metricRepository != nil {
		return s.metricRepository
	}
	s.metricRepository = &MetricRepository{
		metrics:    make(map[string]*model.Metric),
		gaugeMap:   GaugeMap(make(map[string]Gauge)),
		counterMap: CounterMap(make(map[string]Counter)),
	}

	return s.metricRepository
}

// s.Metric().Save()

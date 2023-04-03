package teststore

import (
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
		gaugeMap:   GaugeMap(make(map[string]float64)),
		counterMap: CounterMap(make(map[string]int64)),
	}

	return s.metricRepository
}

// s.Metric().Save()

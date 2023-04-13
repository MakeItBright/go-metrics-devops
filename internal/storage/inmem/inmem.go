package inmem

import (
	"context"
	"errors"
	"path"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

type MetricPath string

func Path(mt model.MetricType, mn model.MetricName) MetricPath {
	return MetricPath(path.Join(string(mt), string(mn)))
}

// Store ...
type Store struct {
	// metricRepository *MetricRepository
	db map[MetricPath]model.Metric
}

func New() *Store {
	return &Store{}
}

func (s *Store) MetricStore(_ context.Context, m model.Metric) error {
	s.db[Path(m.Type, m.Name)] = m
	return nil
}
func (s *Store) MetricFetch(_ context.Context, mt model.MetricType, mn model.MetricName) (model.Metric, error) {
	m, ok := s.db[Path(mt, mn)]
	if !ok {
		return model.Metric{}, errors.New("error db")
	}
	return m, nil
}

// Metric ...
// func (s *Store) Metric() store.MetricRepository {

// 	if s.metricRepository != nil {
// 		return s.metricRepository
// 	}
// 	s.metricRepository = &MetricRepository{
// 		gaugeMap:   GaugeMap(make(map[string]float64)),
// 		counterMap: CounterMap(make(map[string]int64)),
// 	}

// 	return s.metricRepository
// }

// s.Metric().Save()

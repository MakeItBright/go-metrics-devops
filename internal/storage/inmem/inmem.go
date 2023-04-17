package inmem

import (
	"context"
	"errors"
	"path"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

func Path(mt model.MetricType, mn model.MetricName) model.MetricPath {
	return model.MetricPath(path.Join(string(mt), string(mn)))
}

// Store ...
type Store struct {
	// metricRepository *MetricRepository
	db map[model.MetricPath]model.Metric
}

func New() *Store {
	return &Store{
		db: make(map[model.MetricPath]model.Metric),
	}
}

func (s *Store) MetricStore(_ context.Context, m model.Metric) error {
	s.db[Path(m.Type, m.Name)] = m
	return nil
}

func (s *Store) MetricFetch(
	_ context.Context,
	mt model.MetricType,
	mn model.MetricName,
) (model.Metric, error) {
	m, ok := s.db[Path(mt, mn)]
	if !ok {
		return model.Metric{}, errors.New("error db")
	}
	return m, nil
}

func (s *Store) MetricAll() map[model.MetricPath]model.Metric {
	return s.db
}

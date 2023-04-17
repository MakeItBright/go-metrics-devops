package storage

import (
	"path"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

// Storage - интерфейс для хранения метрик
type Storage interface {
	AddGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAllMetrics() map[model.MetricPath]model.Metric
}

// MemStorage - реализация интерфейса Storage в памяти.
type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
}

// NewMemStorage создает новый экземпляр MemStorage.
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gaugeMap:   make(map[string]float64),
		counterMap: make(map[string]int64),
	}
}

// AddGauge добавляет новое значение в map для метрики типа gauge.
func (s *MemStorage) AddGauge(name string, value float64) {
	s.gaugeMap[name] = value
}

// AddCounter добавляет новое значение в map для метрики типа counter.
func (s *MemStorage) AddCounter(name string, value int64) {
	s.counterMap[name] += value
}

// GetGauge возвращает значение метрики типа gauge из map.
// Второе возвращаемое значение - флаг наличия метрики в map.
func (s *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := s.gaugeMap[name]
	return val, ok
}

// GetCounter возвращает значение метрики типа counter из map.
// Второе возвращаемое значение - флаг наличия метрики в map.
func (s *MemStorage) GetCounter(name string) (int64, bool) {
	val, ok := s.counterMap[name]
	return val, ok
}

// GetAllMetrics возвращает все метрики в map.
func (s *MemStorage) GetAllMetrics() map[model.MetricPath]model.Metric {

	metrics := make(map[model.MetricPath]model.Metric)

	for name, val := range s.gaugeMap {
		metrics[Path(model.MetricTypeGauge, model.MetricName(name))] = model.Metric{
			Name:  model.MetricName(name),
			Type:  model.MetricTypeGauge,
			Value: val,
		}
	}

	for name, val := range s.counterMap {
		metrics[Path(model.MetricTypeCounter, model.MetricName(name))] = model.Metric{
			Name:  model.MetricName(name),
			Type:  model.MetricTypeCounter,
			Delta: val,
		}
	}

	return metrics
}

func Path(mt model.MetricType, mn model.MetricName) model.MetricPath {
	return model.MetricPath(path.Join(string(mt), string(mn)))
}

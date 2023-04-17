package monitor

import (
	"math/rand"
	"runtime"

	"github.com/MakeItBright/go-metrics-devops/internal/sender"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/pkg/errors"
)

// Monitor представляет собой структуру для сбора и отправки метрик.
type Monitor struct {
	storage storage.Storage
	sender  sender.Sender
}

// NewMonitor создает новый экземпляр Monitor.
func NewMonitor(storage storage.Storage, sender sender.Sender) *Monitor {
	return &Monitor{
		storage: storage,
		sender:  sender,
	}
}

// CollectMetrics собирает метрики и сохраняет их в хранилище.
func (m *Monitor) CollectMetrics() error {
	if err := m.collectRuntimeMetrics(); err != nil {
		return errors.Wrap(err, "failed to collect runtime metrics")
	}

	if err := m.collectSystemMetrics(); err != nil {
		return errors.Wrap(err, "failed to collect system metrics")
	}

	return nil
}

// collectRuntimeMetrics собирает метрики, связанные с работой приложения и сохраняет их в хранилище.
func (m *Monitor) collectRuntimeMetrics() error {
	// здесь логика сбора метрик и сохранения их в storage
	// get runtime metrics
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	// memory metrics
	m.storage.AddGauge("Alloc", float64(mem.Alloc))
	m.storage.AddGauge("TotalAlloc", float64(mem.TotalAlloc))
	m.storage.AddGauge("Sys", float64(mem.Sys))
	m.storage.AddGauge("Lookups", float64(mem.Lookups))
	m.storage.AddGauge("Mallocs", float64(mem.Mallocs))
	m.storage.AddGauge("Frees", float64(mem.Frees))

	// heap memory metrics
	m.storage.AddGauge("HeapAlloc", float64(mem.HeapAlloc))
	m.storage.AddGauge("HeapSys", float64(mem.HeapSys))
	m.storage.AddGauge("HeapIdle", float64(mem.HeapIdle))
	m.storage.AddGauge("HeapInuse", float64(mem.HeapInuse))
	m.storage.AddGauge("HeapReleased", float64(mem.HeapReleased))
	m.storage.AddGauge("HeapObjects", float64(mem.HeapObjects))

	// stack memory metrics
	m.storage.AddGauge("StackInuse", float64(mem.StackInuse))
	m.storage.AddGauge("StackSys", float64(mem.StackSys))

	// GC metrics
	m.storage.AddGauge("NumGC", float64(mem.NumGC))
	m.storage.AddGauge("PauseTotalNs", float64(mem.PauseTotalNs))
	m.storage.AddGauge("LastGC", float64(mem.LastGC))
	m.storage.AddGauge("NextGC", float64(mem.NextGC))

	m.storage.AddGauge("RandomValue", rand.Float64())

	return nil
}

// collectSystemMetrics собирает метрики, связанные с системными ресурсами и сохраняет их в хранилище.
func (m *Monitor) collectSystemMetrics() error {
	// здесь логика сбора метрик и сохранения их в storage
	// get system metrics like random and counter

	m.storage.AddCounter("PollCounter", 1)

	return nil
}

// Dump получает все метрики из хранилища и отправляет их на сервер.
func (m *Monitor) Dump() error {
	metrics := m.storage.GetAllMetrics()

	if err := m.sender.SendMetrics(metrics); err != nil {
		return errors.Wrap(err, "failed to send metrics")
	}

	return nil
}

package agent

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/sender"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
)

// agent представляет собой структуру для сбора и отправки метрик.
type agent struct {
	storage storage.Storage
	sender  *sender.Sender
}

// Newagent создает новый экземпляр agent.
func NewAgent(storage storage.Storage, sender *sender.Sender) *agent {
	return &agent{
		storage: storage,
		sender:  sender,
	}
}

// Run запуск агента с переданной конфигурацией.
func Start(cfg Config) error {

	// устанавливаем интервал для периодической отправки HTTP-запросов
	pollTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	defer pollTicker.Stop()

	// устанавливаем интервал для отправки метрик
	reportTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
	defer reportTicker.Stop()

	a := NewAgent(
		storage.NewMemStorage(),
		sender.NewSender(cfg.Scheme+"://"+cfg.Address),
	)

	// запускаем бесконечный цикл для периодической отправки HTTP-запросов
	for {

		select {
		case <-pollTicker.C:
			log.Printf("agent is running, collect metrics every %v seconds", cfg.PollInterval)
			a.CollectMetrics()

		case <-reportTicker.C:
			log.Printf("agent is running, sending requests to %v every %v seconds", cfg.Address, cfg.ReportInterval)
			if err := a.Dump(); err != nil {
				log.Printf("ERROR: cannot agent dump: %s", err)
			}

		}
	}
}

// CollectMetrics собирает метрики и сохраняет их в хранилище.
func (a *agent) CollectMetrics() {
	a.collectRuntimeMetrics()
	a.collectSystemMetrics()

}

// collectRuntimeMetrics собирает метрики, связанные с работой приложения и сохраняет их в хранилище.
func (a *agent) collectRuntimeMetrics() {

	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	// memory metrics
	a.storage.AddGauge("Alloc", float64(mem.Alloc))
	a.storage.AddGauge("TotalAlloc", float64(mem.TotalAlloc))
	a.storage.AddGauge("Sys", float64(mem.Sys))
	a.storage.AddGauge("Lookups", math.Max(float64(mem.Lookups), 1))
	a.storage.AddGauge("Mallocs", float64(mem.Mallocs))
	a.storage.AddGauge("Frees", float64(mem.Frees))

	// heap memory metrics
	a.storage.AddGauge("HeapAlloc", float64(mem.HeapAlloc))
	a.storage.AddGauge("HeapSys", float64(mem.HeapSys))
	a.storage.AddGauge("HeapIdle", float64(mem.HeapIdle))
	a.storage.AddGauge("HeapInuse", float64(mem.HeapInuse))
	a.storage.AddGauge("HeapReleased", float64(mem.HeapReleased))
	a.storage.AddGauge("HeapObjects", float64(mem.HeapObjects))

	// stack memory metrics
	a.storage.AddGauge("StackInuse", float64(mem.StackInuse))
	a.storage.AddGauge("StackSys", float64(mem.StackSys))

	// GC metrics
	a.storage.AddGauge("NumGC", float64(mem.NumGC))
	a.storage.AddGauge("PauseTotalNs", float64(mem.PauseTotalNs))
	a.storage.AddGauge("LastGC", float64(mem.LastGC))
	a.storage.AddGauge("NextGC", float64(mem.NextGC))

	a.storage.AddGauge("BuckHashSys", float64(mem.BuckHashSys))
	a.storage.AddGauge("GCCPUFraction", float64(mem.GCCPUFraction))
	a.storage.AddGauge("GCSys", float64(mem.GCSys))

	a.storage.AddGauge("MCacheInuse", float64(mem.MCacheInuse))
	a.storage.AddGauge("MCacheSys", float64(mem.MCacheSys))
	a.storage.AddGauge("MSpanInuse", float64(mem.MSpanInuse))
	a.storage.AddGauge("MSpanSys", float64(mem.MSpanSys))
	a.storage.AddGauge("NumForcedGC", float64(mem.NumForcedGC))
	a.storage.AddGauge("OtherSys", float64(mem.OtherSys))

	// TODO remove this hacks
	runtime.GC()

	a.storage.AddGauge("RandomValue", rand.Float64())

}

// collectSystemMetrics собирает метрики, связанные с системными ресурсами и сохраняет их в хранилище.
func (a *agent) collectSystemMetrics() {
	// здесь логика сбора метрик и сохранения их в storage
	// get system metrics like random and counter

	a.storage.AddCounter("PollCounter", 1)
	a.storage.AddCounter("PollCount", 1)

}

// Dump получает все метрики из хранилища и отправляет их на сервер.
func (a *agent) Dump() error {
	metrics := a.storage.GetAllMetrics()

	if err := a.sender.SendMetrics(metrics); err != nil {
		return fmt.Errorf("cannot send metrics: %w", err)
	}

	return nil

}

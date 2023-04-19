package server

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// Server ...
type server struct {
	logger *logrus.Logger
	router *chi.Mux
	sm     storage.Storage
}

// Metric ...
type Metric struct {
	Name  string  `json:"id"`              // имя метрики
	Type  string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// New ...
func newServer(sm storage.Storage) *server {
	s := &server{
		logger: logrus.New(),
		router: chi.NewRouter(),
		sm:     sm,
	}
	s.configureRouter()

	return s
}

// ServeHTTP
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Router return chi.Router for testing and actual work
func (s *server) configureRouter() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.StripSlashes)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/", s.handleGetAllMetrics)
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.handlePostUpdateMetric)
	s.router.Get("/value/{metricType}/{metricName}", s.handleGetMetric)
}

// handlePostUpdateMetric
func (s *server) handlePostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	switch metricType {
	case string(model.MetricTypeGauge):
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.sm.AddGauge(metricName, value)

	case string(model.MetricTypeCounter):
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.sm.AddCounter(metricName, delta)

	default:

		w.WriteHeader(http.StatusBadRequest)
		return

	}

	// response answer
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`Metric updated`))

}

// handleGetAllMetrics  возвращающая все имеющиеся метрики и их значения в виде HTML-страницы
func (s *server) handleGetAllMetrics(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("All Metrics")
	tmpl, err := template.ParseFiles("templates/index.go.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, s.sm.GetAllMetrics())
}

// handleGetMetricво значение метрики
func (s *server) handleGetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	w.Header().Set("Content-Type", "text/plain")

	switch metricType {
	case string(model.MetricTypeGauge):
		value, ok := s.sm.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", value)))

	case string(model.MetricTypeCounter):
		delta, ok := s.sm.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", delta)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

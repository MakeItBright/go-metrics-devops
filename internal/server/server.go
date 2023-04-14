package server

import (
	"fmt"
	"io"
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
	sm     storage.Metric
}

// Metric ...
type Metric struct {
	Name  string  `json:"id"`              // имя метрики
	Type  string  `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// New ...
func newServer(sm storage.Metric) *server {
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
	s.router.Get("/health", s.handleHealth())
	s.router.Get("/", s.handleGetAllMetrics())
	s.router.Post("/update/{metricType}/{metricName}/{metricValue}", s.handlePostUpdateMetric())
	s.router.Get("/value/{metricType}/{metricName}", s.handleGetMetric())
}

// handlePostUpdateMetric
func (s *server) handlePostUpdateMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		metricValue := chi.URLParam(r, "metricValue")

		if metricType != "gauge" && metricType != "counter" {
			http.Error(w, "Не поддерживаемый тип метрики", http.StatusNotImplemented)
			return
		}

		var (
			mt    model.MetricType
			delta int64   // значение метрики в случае передачи counter
			value float64 // значение метрики в случае передачи gauge
			err   error
		)

		switch metricType {
		case "gauge":
			mt = model.MetricTypeGauge
			value, err = strconv.ParseFloat(metricValue, 64)

		case "counter":
			mt = model.MetricTypeCounter
			delta, err = strconv.ParseInt(metricValue, 10, 64)

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// http.Error(w, err.Error(), 500)
			return
		}
		m, err := s.sm.MetricFetch(r.Context(), model.MetricType(metricType), model.MetricName(metricName))
		if err != nil {
			if err = s.sm.MetricStore(r.Context(), model.Metric{
				Name:  model.MetricName(metricName),
				Type:  mt,
				Delta: delta,
				Value: value,
			}); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				// http.Error(w, err.Error(), 500)
				return
			}
		} else {
			if err = s.sm.MetricStore(r.Context(), model.Metric{
				Name:  model.MetricName(metricName),
				Type:  mt,
				Delta: delta + m.Delta,
				Value: value,
			}); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				// http.Error(w, err.Error(), 500)
				return
			}
		}

		// response answer
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`Metric updated`))
	}
}

// handleGetMetricво значение метрики
func (s *server) handleGetMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Get Metric")
		metricType := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "metricName")
		m, err := s.sm.MetricFetch(r.Context(), model.MetricType(metricType), model.MetricName(metricName))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")

		if metricType != "gauge" && metricType != "counter" {
			http.Error(w, "Не поддерживаемый тип метрики", http.StatusNotImplemented)
			return
		}
		switch metricType {
		case "gauge":
			w.Write([]byte(fmt.Sprintf("%v", m.Value)))
		case "counter":
			w.Write([]byte(fmt.Sprintf("%d", m.Delta)))
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}

// handleGetAllMetrics  возвращающая все имеющиеся метрики и их значения в виде HTML-страницы
func (s *server) handleGetAllMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("All Metrics")
		tmpl, err := template.ParseFiles("templates/index.go.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		s.logger.Info(s.sm.MetricAll())
		tmpl.Execute(w, s.sm.MetricAll())
	}
}

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Test Health")
		io.WriteString(w, "Test Health")

	}
}

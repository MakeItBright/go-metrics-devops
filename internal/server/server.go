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

type server struct {
	logger *logrus.Logger
	router *chi.Mux
	sm     storage.Storage
}

func newServer(sm storage.Storage) *server {
	s := &server{
		logger: logrus.New(),
		router: chi.NewRouter(),
		sm:     sm,
	}

	s.registerRouter()

	return s
}

// ServeHTTP
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Router return chi.Router for testing and actual work
func (s *server) registerRouter() {
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

	switch model.MetricType(metricType) {
	case model.MetricTypeGauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			s.logger.Errorf("cannot parse gauge metric value: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.sm.AddGauge(metricName, value)

	case model.MetricTypeCounter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			s.logger.Errorf("cannot parse counter metric value: %s", err)
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
	tmpl, err := template.ParseFiles("templates/index.go.html")
	if err != nil {
		s.logger.Errorf("cannot parse template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, s.sm.GetAllMetrics()); err != nil {
		s.logger.Errorf("cannot execute template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// handleGetMetricво значение метрики
func (s *server) handleGetMetric(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, "metricType")
	metricName := chi.URLParam(r, "metricName")

	w.Header().Set("Content-Type", "text/plain")

	switch model.MetricType(metricType) {
	case model.MetricTypeGauge:
		value, ok := s.sm.GetGauge(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", value)))

	case model.MetricTypeCounter:
		delta, ok := s.sm.GetCounter(metricName)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%v", delta)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

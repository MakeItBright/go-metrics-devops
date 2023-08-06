package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	logger "github.com/MakeItBright/go-metrics-devops/internal/logger"
	mw "github.com/MakeItBright/go-metrics-devops/internal/middleware"
	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type server struct {
	router *chi.Mux
	sm     storage.Storage
}

func newServer(sm storage.Storage) *server {
	s := &server{
		router: chi.NewRouter(),
		sm:     sm,
	}

	s.registerRouter()

	return s
}

// ServeHTTP реализует интерфейс http.Handler и обрабатывает HTTP-запросы
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}

// Router return chi.Router for testing and actual work
func (s *server) registerRouter() {

	s.router.Use(logger.RequestLogger)
	s.router.Use(middleware.StripSlashes)
	s.router.Use(mw.GzipMiddleware)
	s.router.Get("/health", s.handleHealth)
	s.router.Get("/", s.handleGetAllMetrics)

	s.router.Post("/update", s.handleJSONPostUpdateMetric)
	s.router.Post("/value", s.handleJSONPostGetMetric)

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
			logger.Log.Sugar().Errorf("cannot parse gauge metric value: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.sm.AddGauge(metricName, value)

	case model.MetricTypeCounter:
		delta, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			logger.Log.Sugar().Errorf("cannot parse counter metric value: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		s.sm.AddCounter(metricName, delta)

	default:

		w.WriteHeader(http.StatusNotImplemented)
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
		logger.Log.Sugar().Errorf("cannot parse template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, s.sm.GetAllMetrics()); err != nil {
		logger.Log.Sugar().Errorf("cannot execute template: %s", err)
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

// handleJSONPostUpdateMetric
func (s *server) handleJSONPostUpdateMetric(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		logger.Log.Sugar().Info("wrong content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Log.Sugar().Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var m model.Metric

	err = json.Unmarshal(body, &m)

	if err != nil {
		logger.Log.Sugar().Errorf("cannot parse metric: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch model.MetricType(m.Type) {
	case model.MetricTypeGauge:
		s.sm.AddGauge(string(m.Name), m.Value)
		value, ok := s.sm.GetGauge(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m.Value = value

	case model.MetricTypeCounter:
		s.sm.AddCounter(string(m.Name), m.Delta)
		delta, ok := s.sm.GetCounter(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		m.Delta = delta

	default:

		w.WriteHeader(http.StatusBadRequest)
		return

	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// handleGetMetric
func (s *server) handleJSONPostGetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Log.Sugar().Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var m model.Metric

	err = json.Unmarshal(body, &m)
	if err != nil {
		logger.Log.Sugar().Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if m.Name == "" {
		logger.Log.Sugar().Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch model.MetricType(m.Type) {
	case model.MetricTypeGauge:
		value, ok := s.sm.GetGauge(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Value = value

	case model.MetricTypeCounter:
		delta, ok := s.sm.GetCounter(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Delta = delta

	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

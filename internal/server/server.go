package server

import (
	"encoding/json"
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

// ServeHTTP реализует интерфейс http.Handler и обрабатывает HTTP-запросы
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)

}

// Router return chi.Router for testing and actual work
func (s *server) registerRouter() {

	// Создание и конфигурирование логгера
	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.InfoLevel

	// Middleware для логирования
	// s.router.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		start := time.Now()

	// 		// Создание обертки над ResponseWriter, чтобы сохранить данные о коде статуса и размере содержимого
	// 		rw := newResponseWriter(w)

	// 		// Выполнение запроса
	// 		next.ServeHTTP(rw, r)

	// 		// Запись лога
	// 		logger.WithFields(logrus.Fields{
	// 			"method":         r.Method,
	// 			"uri":            r.RequestURI,
	// 			"elapsed_ms":     time.Since(start).Microseconds(),
	// 			"status_code":    rw.statusCode,
	// 			"response_bytes": rw.contentLength,
	// 		}).Info("Request processed")
	// 	})
	// })

	// s.router.Use(WithLogging)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.StripSlashes)
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

// handleJSONPostUpdateMetric
func (s *server) handleJSONPostUpdateMetric(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("||| ================= POST Update ======================== ||||")
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		s.logger.Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var m model.Metric

	err = json.Unmarshal(body, &m)
	s.logger.Printf("body to m: %+v", m)
	if err != nil {
		s.logger.Errorf("cannot parse metric: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch model.MetricType(m.Type) {
	case model.MetricTypeGauge:
		// value, err := strconv.ParseFloat(m.Value, 64)
		// if err != nil {
		// 	s.logger.Errorf("cannot parse gauge metric value: %s", err)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		s.sm.AddGauge(string(m.Name), m.Value)
		value, ok := s.sm.GetGauge(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Value = value
	case model.MetricTypeCounter:
		// delta, err := strconv.ParseInt(m.Delta, 10, 64)
		// if err != nil {
		// 	s.logger.Errorf("cannot parse counter metric value: %s", err)
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }

		s.sm.AddCounter(string(m.Name), m.Delta)
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		s.logger.Errorf("error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// handleGetMetricво значение метрики
func (s *server) handleJSONPostGetMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.logger.Info("||| ================= POST Value ======================== ||||")
	body, err := io.ReadAll(r.Body)

	if err != nil {
		s.logger.Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var m model.Metric

	err = json.Unmarshal(body, &m)
	s.logger.Info(m)
	s.logger.Info(err)
	if err != nil {
		s.logger.Errorf("cannot parse counter metric value: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if m.Name == "" {
		s.logger.Errorf("cannot parse counter metric value: %s", err)
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
		// w.Write([]byte(fmt.Sprintf("%v", value)))

	case model.MetricTypeCounter:
		delta, ok := s.sm.GetCounter(string(m.Name))
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		m.Delta = delta
		// w.Write([]byte(fmt.Sprintf("%v", delta)))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		s.logger.Errorf("error: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

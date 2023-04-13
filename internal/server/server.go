package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server ...
type server struct {
	logger *logrus.Logger
	router *mux.Router
	sm     storage.Metric
}

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
		router: mux.NewRouter(),
		sm:     sm,
	}
	s.configureRouter()

	return s
}

// ServeHTTP
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Config Router ...
func (s *server) configureRouter() {
	s.router.HandleFunc("/health", s.handleHealth())
	/// update/counter/someMetric/527 HTTP/1.1
	s.router.HandleFunc("/update/{mtype}/{mname}/{mvalue}", s.handlePostUpdateMetric()).Methods("POST")

}
func (s *server) handlePostUpdateMetric() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		s.logger.Info("Update")
		s.logger.Info(params)
		var (
			mt    model.MetricType
			delta int64   // значение метрики в случае передачи counter
			value float64 // значение метрики в случае передачи gauge
			err   error
		)

		switch params["mtype"] {
		case "gauge":
			mt = model.MetricTypeGauge
			delta, err = strconv.ParseInt(params["mvalue"], 10, 64)

		case "counter":
			mt = model.MetricTypeCounter
			value, err = strconv.ParseFloat(params["mvalue"], 64)

		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// http.Error(w, err.Error(), 500)
			return
		}

		if err = s.sm.MetricStore(r.Context(), model.Metric{
			Name:  model.MetricName(params["mname"]),
			Type:  mt,
			Delta: delta,
			Value: value,
		}); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			// http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Test Health")

		io.WriteString(w, "Test Health")

	}
}

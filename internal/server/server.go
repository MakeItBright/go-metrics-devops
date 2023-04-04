package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/MakeItBright/go-metrics-devops/internal/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server ...
type server struct {
	logger *logrus.Logger
	router *mux.Router
	store  store.Store
}

// New ...
func newServer(store store.Store) *server {
	s := &server{
		logger: logrus.New(),
		router: mux.NewRouter(),
		store:  store,
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
	s.router.HandleFunc("/update/{mtype}/{mname}/{mvalue}", s.handlePostUpdateMetric()).Methods("GET", "POST")

}

func (s *server) handlePostUpdateMetric() http.HandlerFunc {

	// type request struct {
	// 	Name  string // имя метрики
	// 	MType string // параметр, принимающий значение gauge или counter
	// }
	type Metric struct {
		MName string  // имя метрики
		MType string  // параметр, принимающий значение gauge или counter
		Delta int64   // значение метрики в случае передачи counter
		Value float64 // значение метрики в случае передачи gauge
	}
	return func(w http.ResponseWriter, r *http.Request) {

		s.logger.Info("HandlePostUpdateMetric")

		s.logger.Info(mux.Vars(r))
		var metrics Metric
		// Convert the map to JSON
		jsonData, _ := json.Marshal(mux.Vars(r))
		json.Unmarshal(jsonData, &metrics)

		s.logger.Info(metrics)
		s.logger.Info(metrics.MName)
		s.logger.Info(metrics.MType)
		s.logger.Info(metrics.Value)

		switch metrics.MType {
		case "Gauge":
		case "Counter":
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// if err != nil {
		// 	http.Error(w, err.Error(), 500)
		// 	return
		// }

		// store.MetricRepository.SaveGaugeValue()
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		//ttp.StatusNotFound
	}
}

func (s *server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("Test Health")

		io.WriteString(w, "Test Health")
	}
}

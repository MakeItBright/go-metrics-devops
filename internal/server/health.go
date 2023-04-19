package server

import (
	"io"
	"net/http"
)

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Test Health")
	io.WriteString(w, "Test Health")

}

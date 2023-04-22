package server

import (
	"io"
	"net/http"
)

func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Ok")

}

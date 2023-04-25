package server

import "net/http"

// Обертка над ResponseWriter, чтобы сохранить данные о коде статуса и размере содержимого
type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK, 0}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	length, err := rw.ResponseWriter.Write(b)
	rw.contentLength += length
	return length, err
}

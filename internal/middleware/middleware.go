package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

var ContentTypesForCompress = map[string]bool{
	"application/json": true,
	"text/html":        true,
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Encoding") == "gzip" {
		return w.Writer.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w gzipWriter) WriteHeader(code int) {
	contentType := w.Header().Get("Content-Type")
	enableCompress := ContentTypesForCompress[contentType]

	if enableCompress {
		w.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(code)
}

func GZipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// проверяем что запрос пришел сжатый
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(gz)
			defer gz.Close()
		}

		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем типы контента, для которых применяется сжатие
		// contentType := w.Header().Get("Content-Type")
		// enableCompress := ContentTypesForCompress[contentType]

		// if !enableCompress {
		// 	// Если тип контента не соответствует, передаем управление
		// 	// дальше без изменений
		// 	next.ServeHTTP(w, r)
		// 	return
		// }
		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		// устанавливаем соответствующие заголовки сервера
		w.Header().Set("Content-Encoding", "gzip")

		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}

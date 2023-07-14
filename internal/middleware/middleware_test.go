package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGZipHandle(t *testing.T) {
	// Создаем новый тестовый HTTP-сервер
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что контент был сжат
		if r.Header.Get("Content-Encoding") != "gzip" {
			t.Errorf("Expected Content-Encoding to be gzip, got %s", r.Header.Get("Content-Encoding"))
		}

		// Проверяем, что типы контента соответствуют application/json и text/html
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" && contentType != "text/html" {
			t.Errorf("Expected Content-Type to be application/json or text/html, got %s", contentType)
		}

		// Читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}

		// Проверяем, что содержимое тела запроса правильное
		expectedBody := "Test Request Body"
		if string(body) != expectedBody {
			t.Errorf("Expected request body to be %s, got %s", expectedBody, string(body))
		}

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewServer(GZipHandle(handler))
	defer server.Close()

	// Создаем запрос с сжатым телом
	body := "Test Request Body"
	gzippedBody := compressData(body)

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(gzippedBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")

	// Отправляем запрос на сервер
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что получили успешный ответ
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

// compressData сжимает данные с использованием gzip
func compressData(data string) []byte {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(data)); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

// Writer представляет собой поставщик метрик.
type Writer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewWriter создает новый экземпляр Writer и инициализирует его.
// Принимает fileName - имя файла, в который будут сохраняться метрики.
// Возвращает указатель на созданный Writer и ошибку, если возникла.
func NewWriter(fileName string) (*Writer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Writer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteMetrics записывает переданные метрики в файл.
// Принимает metrics - срез метрик для записи.
func (p *Writer) WriteMetrics(metrics []model.Metric) error {
	if metrics != nil {
		return p.encoder.Encode(metrics)
	}
	return errors.New("can't write metric to file from memory - object is empty")
}

// Close закрывает файловый дескриптор и освобождает ресурсы.
// Возвращает ошибку, если возникла.
func (p *Writer) Close() error {
	return p.file.Close()
}

// Reader представляет собой структуру, ответственную за чтение метрик из файла.
type Reader struct {
	file    *os.File
	decoder *json.Decoder
}

// NewReader создает новый экземпляр Reader и открывает файл для чтения.
// fileName - полное имя файла, из которого будут читаться метрики.
// Возвращает указатель на Reader и ошибку, если возникла.
func NewReader(fileName string) (*Reader, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadMetrics считывает метрики из файла и возвращает их в виде среза моделей Metric.
// Возвращает срез метрик и ошибку, если возникла.
func (c *Reader) ReadMetrics() ([]model.Metric, error) {
	metrics := []model.Metric{}

	if err := c.decoder.Decode(&metrics); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}
	return metrics, nil
}

// Close закрывает файл после чтения метрик.
// Возвращает ошибку, если возникла при закрытии файла.
func (c *Reader) Close() error {
	return c.file.Close()
}

package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

// Producer представляет собой поставщик метрик.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer создает новый экземпляр Producer и инициализирует его.
// Принимает fileName - имя файла, в который будут сохраняться метрики.
// Возвращает указатель на созданный Producer и ошибку, если возникла.
func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

// WriteMetrics записывает переданные метрики в файл.
// Принимает metrics - срез метрик для записи.
func (p *Producer) WriteMetrics(metrics []model.Metric) error {
	if metrics != nil {
		return p.encoder.Encode(metrics)
	}
	return errors.New("can't write metric to file from memory - object is empty")
}

// Close закрывает файловый дескриптор и освобождает ресурсы.
// Возвращает ошибку, если возникла.
func (p *Producer) Close() error {
	return p.file.Close()
}

// consumer представляет собой структуру, ответственную за чтение метрик из файла.
type consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer создает новый экземпляр consumer и открывает файл для чтения.
// fileName - полное имя файла, из которого будут читаться метрики.
// Возвращает указатель на consumer и ошибку, если возникла.
func NewConsumer(fileName string) (*consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

// ReadMetrics считывает метрики из файла и возвращает их в виде среза моделей Metric.
// Возвращает срез метрик и ошибку, если возникла.
func (c *consumer) ReadMetrics() ([]model.Metric, error) {
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
func (c *consumer) Close() error {
	return c.file.Close()
}

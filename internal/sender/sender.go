package sender

import (
	"fmt"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/go-resty/resty/v2"
)

type Sender struct {
	url string
}

// NewSender создает новый экземпляр Sender.
func NewSender(serverAddr string) *Sender {
	return &Sender{
		url: serverAddr,
	}
}

// SendMetrics отправляет метрики на сервер.
func (s *Sender) SendMetrics(metrics []model.Metric) error {
	client := resty.New()
	client.
		// устанавливаем количество повторений
		SetRetryCount(3).
		// длительность ожидания между попытками
		SetRetryWaitTime(30 * time.Second).
		// длительность максимального ожидания
		SetRetryMaxWaitTime(90 * time.Second)
	for _, value := range metrics {
		url := fmt.Sprintf("%s/update", s.url)
		_, err := client.R().SetHeader("Content-Type", "application/json").SetBody(value).Post(url)

		if err != nil {
			return fmt.Errorf("cannot perform POST request: %w", err)
		}

	}

	return nil
}

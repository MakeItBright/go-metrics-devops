package sender

import (
	"fmt"

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
	for _, value := range metrics {
		url := fmt.Sprintf("%s/update/%s/%s/%s", s.url, value.Type, value.Name, value.GetValue())

		_, err := client.R().SetHeader("Content-Type", "text/plain").Post(url)
		if err != nil {
			return fmt.Errorf("cannot perform POST request: %w", err)
		}

	}
	return nil
}

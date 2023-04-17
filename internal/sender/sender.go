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
func (s *Sender) SendMetrics(metrics map[model.MetricPath]model.Metric) error {
	client := resty.New()
	for _, value := range metrics {
		var url string
		switch value.Type {
		case "gauge":
			url = fmt.Sprintf("%s/update/%s/%s/%v", s.url, value.Type, value.Name, value.Value)
		case "counter":
			url = fmt.Sprintf("%s/update/%s/%s/%v", s.url, value.Type, value.Name, value.Delta)
		default:

		}
		_, err := client.R().
			SetHeader("Content-Type", "text/plain").
			Post(url)
		if err != nil {
			return err
		}
		fmt.Println(url)
	}
	return nil
}

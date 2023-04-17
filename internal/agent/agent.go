package agent

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Run запуск агента с переданной конфигурацией.
func Run(cfg Config) error {
	if cfg.Logger == nil {
		cfg.Logger = logrus.StandardLogger()
	}

	// создаем REST-клиент для отправки HTTP-запросов
	client := resty.New()
	urls := make([]string, 29)
	host := "http://" + cfg.Address

	// устанавливаем интервал для периодической отправки HTTP-запросов
	pollInterval := cfg.PollInterval
	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	// устанавливаем интервал для отправки метрик
	reportInterval := cfg.ReportInterval
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	i := 0
	// запускаем бесконечный цикл для периодической отправки HTTP-запросов
	for {

		select {
		case <-pollTicker.C:
			// собираем метрики
			cfg.Logger.Infof(
				"agent is running, collect metrics every %v seconds",
				pollInterval.Seconds(),
			)

		case <-reportTicker.C:
			// отправляем HTTP-запросы на указанные адреса
			cfg.Logger.Infof(
				"agent is running, sending requests to %v every %v seconds",
				cfg.Address,
				reportInterval.Seconds(),
			)
			i = 0
			urls[i] = fmt.Sprintf(
				"%s/update/counter/%s/%d",
				host, "PollCount", 1,
			)
			i++
			//RandomValue (тип gauge) — обновляемое произвольное значение.
			urls[i] = fmt.Sprintf(
				"%s/update/gauge/%s/%d",
				host, "RandomValue", rand.Intn(1000000),
			)
			i++
			v := reflect.ValueOf(mem)
			tof := v.Type()

			for j := 0; j < v.NumField(); j++ {

				val := 0.0
				if !v.Field(j).CanUint() && !v.Field(j).CanFloat() {
					continue
				} else if !v.Field(j).CanUint() {
					val = v.Field(j).Float()
				} else {
					val = float64(v.Field(j).Uint())
				}
				name := tof.Field(j).Name
				urls[i] = fmt.Sprintf(
					"%s/update/gauge/%s/%f",
					host, name, val,
				)
				i++
				cfg.Logger.Infof("%s = %f", name, val)

			}

			doRequest(urls, client, cfg.Logger)
		}
	}
}

// doRequest - отправка HTTP-запросов на указанные адреса
func doRequest(urls []string, client *resty.Client, logger *logrus.Logger) {
	for _, url := range urls {
		_, err := client.R().
			Post(url)

		if err != nil {
			logger.Errorf("failed to send request to %v: %v", url, err)
		}
	}
}

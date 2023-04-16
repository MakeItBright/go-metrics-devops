package agent

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"

	"github.com/MakeItBright/go-metrics-devops/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

// Agent представляет агент, который будет запускать метрики
type agent struct {
	logger *logrus.Logger // логгер для отслеживания ошибок
}

// newAgent - создает новый экземпляр агента и возвращает его указатель
func newAgent() *agent {
	a := &agent{
		logger: logrus.New(),
	}
	return a
}

// RunAgent - запуск агента с переданной конфигурацией
func RunAgent(cfg *config.AgentConfig) error {
	// создаем новый агент
	a := newAgent()
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
			a.logger.Infof("agent is running, collect metrics every %v seconds", pollInterval.Seconds())

		case <-reportTicker.C:
			// отправляем HTTP-запросы на указанные адреса
			a.logger.Infof("agent is running, sending requests to %v every %v seconds", cfg.Address, reportInterval.Seconds())
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
				a.logger.Infof("%s = %f", name, val)

			}

			a.logger.Info(urls)
			a.doRequest(urls, client)
		}
	}
}

// doRequest - отправка HTTP-запросов на указанные адреса
func (a *agent) doRequest(urls []string, client *resty.Client) {
	for _, url := range urls {
		_, err := client.R().
			Post(url)

		if err != nil {
			a.logger.Errorf("failed to send request to %v: %v", url, err)
		}
	}
}

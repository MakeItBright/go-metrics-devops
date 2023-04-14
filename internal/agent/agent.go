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

// Agent ...
type agent struct {
	logger *logrus.Logger
}

const (
	pollInterval   = 2 * time.Second
	reportInterval = 2 * time.Second
)

// var metrics = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc","HeapIdle", "HeapInuse", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys","MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs","PollCount", "RandomValue", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

// Start
func Start(config *Config) error {

	client := resty.New()
	urls := make([]string, 29)
	host := "http://localhost:8080"
	ag := newAgent()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	startReport := time.Now()
	i := 0
	for {
		i = 0
		start := time.Now()
		time.Sleep(pollInterval - time.Since(start))
		//PollCount (тип counter) — счётчик, увеличивающийся на 1 при каждом обновлении метрики из пакета runtime
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
			ag.logger.Infof("%s = %f", name, val)

		}

		if reportInterval-time.Since(startReport) <= 0 {
			startReport = time.Now()
			ag.doRequest(urls, client)
		}
	}

}

// New ...
func newAgent() *agent {
	a := &agent{
		logger: logrus.New(),
	}
	return a
}

// doRequest ...
func (a *agent) doRequest(urls []string, client *resty.Client) {
	a.logger.Info("Send")
	for _, url := range urls {
		_, err := client.R().
			Post(url)

		if err != nil {
			panic(err)
		}
	}
}

package agent

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

var (
	flagRunAddr        string // неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
	flagPollInterval   int
	flagReportInterval int
)

func init() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением localhost:8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&flagReportInterval, "r", 10, "interval of report")
	flag.IntVar(&flagPollInterval, "p", 2, "interval of poll")
}

// Agent ...
type agent struct {
	logger *logrus.Logger
}

// var metrics = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc","HeapIdle", "HeapInuse", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys","MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs","PollCount", "RandomValue", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

// Start
func Start(config *Config) error {
	flag.Parse()
	var (
		pollInterval   = 2 * time.Second
		reportInterval = 10 * time.Second
	)
	if envBindAddr := os.Getenv("ADDRESS"); envBindAddr != "" {
		flagRunAddr = envBindAddr
	}
	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {

		r, err := strconv.Atoi(envReportInterval)
		if err != nil {
			log.Fatal(err)
		}
		reportInterval = time.Duration(r) * time.Second
	}
	if envpollInterval := os.Getenv("POLL_INTERVAL"); envpollInterval != "" {

		p, err := strconv.Atoi(envpollInterval)
		if err != nil {
			log.Fatal(err)
		}
		pollInterval = time.Duration(p) * time.Second
	}

	client := resty.New()
	urls := make([]string, 29)
	host := "http://" + flagRunAddr
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

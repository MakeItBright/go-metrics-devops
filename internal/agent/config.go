package agent

// Config структура конфигурации агента.
type Config struct {
	Scheme         string
	Address        string // адрес HTTP-сервера
	PollInterval   int    // интервал обновления метрик
	ReportInterval int    // интервал между отправками метрик
}

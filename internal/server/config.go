package server

// Config структура конфигурации сервера.
type Config struct {
	BindAddr        string
	LogLevel        string
	StoreInterval   int
	FileStoragePath string
	Restore         bool
}

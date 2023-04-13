package server

// Config ...
type Config struct {
	BindAddr string
	LogLevel string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

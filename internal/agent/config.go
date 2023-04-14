package agent

// Config ...
type Config struct {
	LogLevel string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		LogLevel: "debug",
	}
}

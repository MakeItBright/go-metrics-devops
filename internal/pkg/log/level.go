package log

type Level string

const (
	LEVEL_PANIC Level = "panic"
	LEVEL_FATAL Level = "fatal"
	LEVEL_ERROR Level = "error"
	LEVEL_WARN  Level = "warn"
	LEVEL_INFO  Level = "info"
	LEVEL_DEBUG Level = "debug"
	LEVEL_TRACE Level = "trace"
)

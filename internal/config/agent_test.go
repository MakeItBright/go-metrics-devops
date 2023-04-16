package config

import (
	"os"
	"testing"
	"time"
)

// запускаем тесты с переменными окружения
func TestNewAgentConfig(t *testing.T) {
	os.Setenv("ADDRESS", "localhost:8081")
	os.Setenv("POLL_INTERVAL", "5s")
	os.Setenv("REPORT_INTERVAL", "20s")
	defer os.Clearenv()

	cfg := NewAgentConfig()

	expectedAddress := "localhost:8081"
	if cfg.Address != expectedAddress {
		t.Errorf("Address is %s, expected %s", cfg.Address, expectedAddress)
	}

	expectedPollInterval := 5 * time.Second
	if cfg.PollInterval != expectedPollInterval {
		t.Errorf("PollInterval is %v, expected %v", cfg.PollInterval, expectedPollInterval)
	}

	expectedReportInterval := 20 * time.Second
	if cfg.ReportInterval != expectedReportInterval {
		t.Errorf("ReportInterval is %v, expected %v", cfg.ReportInterval, expectedReportInterval)
	}
}

// запускаем тесты с флагами командной строки
func TestNewAgentConfigWithDefaultValues(t *testing.T) {
	cfg := NewAgentConfig()

	expectedAddress := defaultAddress
	if cfg.Address != expectedAddress {
		t.Errorf("Address is %s, expected %s", cfg.Address, expectedAddress)
	}

	expectedPollInterval := defaultPollInterval
	if cfg.PollInterval != expectedPollInterval {
		t.Errorf("PollInterval is %v, expected %v", cfg.PollInterval, expectedPollInterval)
	}

	expectedReportInterval := defaultReportInterval
	if cfg.ReportInterval != expectedReportInterval {
		t.Errorf("ReportInterval is %v, expected %v", cfg.ReportInterval, expectedReportInterval)
	}
}

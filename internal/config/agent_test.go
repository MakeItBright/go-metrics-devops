package config

import (
	"testing"
)

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

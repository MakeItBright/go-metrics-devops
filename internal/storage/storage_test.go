package storage_test

import (
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/storage"
)

func TestMemStorage_AddGauge(t *testing.T) {
	tests := []struct {
		name      string
		gaugeName string
		value     float64
	}{
		{
			name:      "add new gauge",
			gaugeName: "test_gauge",
			value:     1.23,
		},
		{
			name:      "update existing gauge",
			gaugeName: "test_gauge",
			value:     4.56,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			storage.AddGauge(tt.gaugeName, tt.value)
			val, ok := storage.GetGauge(tt.gaugeName)
			if !ok {
				t.Errorf("GetGauge() returned false, expected true")
			}
			if val != tt.value {
				t.Errorf("GetGauge() returned %v, expected %v", val, tt.value)
			}
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	tests := []struct {
		name        string
		counterName string
		value       int64
	}{
		{
			name:        "add new counter",
			counterName: "test_counter",
			value:       123,
		},
		{
			name:        "update existing counter",
			counterName: "test_counter",
			value:       456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := storage.NewMemStorage()
			storage.AddCounter(tt.counterName, tt.value)
			val, ok := storage.GetCounter(tt.counterName)
			if !ok {
				t.Errorf("GetCounter() returned false, expected true")
			}
			if val != tt.value {
				t.Errorf("GetCounter() returned %v, expected %v", val, tt.value)
			}
		})
	}
}

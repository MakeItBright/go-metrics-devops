package storage_test

import (
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_AddGauge(t *testing.T) {
	type input struct {
		name  string
		value float64
	}

	tests := []struct {
		name     string
		input    input
		expected map[string]float64
	}{
		{
			name: "add_gauge_one_metric",
			input: input{
				name:  "metric_name",
				value: 10.0,
			},
			expected: map[string]float64{
				"metric_name": 10.0,
			},
		},
		{
			name: "add_gauge_multiple_metrics",
			input: input{
				name:  "metric_name_2",
				value: 20.0,
			},
			expected: map[string]float64{
				"metric_name":   10.0,
				"metric_name_2": 20.0,
			},
		},
	}

	storage := storage.NewMemStorage()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage.AddGauge(tc.input.name, tc.input.value)

			assert.Equal(t, tc.expected, storage.GetAllMetrics())
		})
	}
}

func TestMemStorage_AddCounter(t *testing.T) {
	type input struct {
		name  string
		value int64
	}

	tests := []struct {
		name     string
		input    input
		expected map[string]interface{}
	}{
		{
			name: "add_counter_one_metric",
			input: input{
				name:  "metric_name",
				value: 1,
			},
			expected: map[string]interface{}{
				"metric_name": int64(1),
			},
		},
		{
			name: "add_counter_multiple_metrics",
			input: input{
				name:  "metric_name_2",
				value: 2,
			},
			expected: map[string]interface{}{
				"metric_name":   int64(1),
				"metric_name_2": int64(2),
			},
		},
		{
			name: "add_counter_same_metrics",
			input: input{
				name:  "metric_name",
				value: 2,
			},
			expected: map[string]interface{}{
				"metric_name":   int64(3),
				"metric_name_2": int64(2),
			},
		},
	}

	storage := storage.NewMemStorage()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage.AddCounter(tc.input.name, tc.input.value)

			assert.Equal(t, tc.expected, storage.GetAllMetrics())
		})
	}
}

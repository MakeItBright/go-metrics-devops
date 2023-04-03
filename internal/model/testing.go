package model

import "testing"

func Testmetric(t *testing.T) *Metric {
	return &Metric{
		Name:  "someMetric",
		MType: "Counter",
		// Delta: 1,
		// Value: 1,
	}
}

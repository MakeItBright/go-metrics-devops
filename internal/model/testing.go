package model

import "testing"

func Testmetric(t *testing.T) *Metric {
	return &Metric{
		Name: "someMetric",
		Type: Counter,
	}
}

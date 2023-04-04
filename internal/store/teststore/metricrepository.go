package teststore

import (
	"errors"
	"fmt"
)

type GaugeMap map[string]float64 // новое значение должно замещать предыдущее.
type CounterMap map[string]int64 // новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.

type MetricRepository struct {
	gaugeMap   GaugeMap
	counterMap CounterMap
}

// SaveCounterValue ...
func (mr *MetricRepository) SaveCounterValue(name string, counter int64) {
	n, ok := mr.counterMap[name]
	if !ok {
		mr.counterMap[name] = counter
		return
	}
	mr.counterMap[name] = n + counter
}

// SaveGaugeValue ...
func (mr *MetricRepository) SaveGaugeValue(name string, gauge float64) {
	mr.gaugeMap[name] = gauge

}

// GetCounterValue ...
func (mr *MetricRepository) GetCounterValue(name string) (int64, error) {
	n, ok := mr.counterMap[name]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return n, nil
}

// GetGaugeValue ...
func (mr *MetricRepository) GetGaugeValue(name string) (float64, error) {
	n, ok := mr.gaugeMap[name]
	if !ok {
		return 0, errors.New("this counter don't find")
	}
	return n, nil
}

// GetAllValues ...
func (mr *MetricRepository) GetAllValues() string {
	mapAll := make(map[string]string)

	for key, val := range mr.counterMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	for key, val := range mr.gaugeMap {
		mapAll[key] = fmt.Sprintf("%v", val)
	}
	var str string
	for key, val := range mapAll {
		str += fmt.Sprintf("%s: %s\n", key, val)
	}
	return str
}

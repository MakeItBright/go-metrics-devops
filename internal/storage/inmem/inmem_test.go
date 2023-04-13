package inmem

import (
	"context"
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
)

func TestStore_MetricStore(t *testing.T) {
	type fields struct {
		db map[MetricPath]model.Metric
	}
	type args struct {
		in0 context.Context
		m   model.Metric
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				db: tt.fields.db,
			}
			if err := s.MetricStore(tt.args.in0, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Store.MetricStore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

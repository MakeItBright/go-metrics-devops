package teststore_test

import (
	"testing"

	"github.com/MakeItBright/go-metrics-devops/internal/model"
	"github.com/MakeItBright/go-metrics-devops/internal/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_UpdateMetric(t *testing.T) {
	ms := teststore.New()

	m := model.Testmetric(t)

	assert.NoError(t, ms.Metric().UpdateMetric(m))
	assert.NotNil(t, m)
}

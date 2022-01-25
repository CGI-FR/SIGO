package sigo_test

import (
	"testing"

	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	t.Parallel()

	values := []float64{21, 23, 5, 24, 15, 23, 19, 24, 7, 10, 21, 22, 22, 23, 24, 23, 23, 24, 25}

	assert.Equal(t, 5.00, sigo.Min(values))
	assert.Equal(t, 25.00, sigo.Max(values))
	assert.Equal(t, 378.00, sigo.Sum(values))
	assert.Equal(t, 19.89, sigo.Mean(values))
	assert.Equal(t, 6.0726614154476515, sigo.Std(values))
}

func TestQuartiles(t *testing.T) {
	t.Parallel()

	values := []float64{21, 23, 5, 24, 15, 23, 19, 24, 7, 10, 21, 22, 22, 23, 24, 23, 23, 24, 25}

	q := sigo.Quartile(values)

	assert.Equal(t, 19.00, q.Q1)
	assert.Equal(t, 24.00, q.Q3)
	assert.Equal(t, 23.00, q.Q2)
	assert.Equal(t, q.Q2, sigo.Median(values))
	assert.Equal(t, 5.00, sigo.IQR(values))
}

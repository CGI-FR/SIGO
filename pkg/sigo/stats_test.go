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

func TestUnique(t *testing.T) {
	t.Parallel()

	values1 := []float64{12, 10, 5, 6, 9, 10, 4, 5, 10, 12, 9, 6, 4, 3, 9, 10}
	values2 := []float64{1, 9, 8, 5, 2, 6, 7, 10, 3, 12, 4, 11}

	res1 := sigo.Unique(values1)
	res2 := sigo.Unique(values2)

	assert.Equal(t, 7, res1)
	assert.Equal(t, 12, res2)
}

func TestOrderMap(t *testing.T) {
	t.Parallel()

	values := map[string]int{"x": 2, "y": 3}

	res := sigo.Order(values)

	assert.Equal(t, []string{"y", "x"}, res)
}

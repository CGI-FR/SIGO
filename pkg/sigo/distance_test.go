package sigo_test

import (
	"math"
	"testing"

	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()

	X := map[string]float64{
		"a": 14, "b": 6, "c": 18, "d": 52.1,
		"e": 21, "f": 36.48, "g": 12, "h": 39,
	}
	Y := map[string]float64{
		"a": 14.17, "b": 6, "c": 4, "d": 12.86,
		"e": 54, "f": 49, "g": 7.2, "h": 27.12,
	}

	dist := sigo.Cosine(X, Y)

	assert.InDelta(t, 0.7622963735959798, dist, math.Pow10(-15))
}

func TestEuclideanDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"k1": 0, "k2": 0}
	Y := map[string]float64{"k1": 2, "k2": 0}

	dist := sigo.Euclidean(X, Y)

	assert.Equal(t, 2.00, dist)
}

func TestManhattanDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 9, "q2": 5}

	dist := sigo.Manhattan(X, Y)

	assert.Equal(t, 10.00, dist)
}

func TestCanberraDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 4, "q2": 4}

	dist := sigo.Camberra(X, Y)

	assert.Equal(t, float64(2)/float64(3), dist)
}

func TestChebyshevDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 9, "q2": 5}

	dist := sigo.Chebyshev(X, Y)

	assert.Equal(t, 7.00, dist)
}

func TestMinkowskiDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 4, "q2": 4}

	dist := sigo.Minkowski(X, Y, 6)

	assert.Equal(t, math.Pow(128, 1.00/6.00), dist)
}

func TestComputeDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"k1": 0, "k2": 0}
	Y := map[string]float64{"k1": 2, "k2": 0}

	dist := sigo.ComputeDistance("", X, Y)

	assert.Equal(t, 2.00, dist)
}

func TestComputeSimilarity(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"k1": 0, "k2": 0}
	Y := map[string]float64{"k1": 2, "k2": 0}

	dist := sigo.ComputeDistance("", X, Y)

	assert.Equal(t, float64(1)/float64(3), sigo.Similarity(dist))
}

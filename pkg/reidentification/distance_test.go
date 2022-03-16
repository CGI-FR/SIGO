package reidentification_test

import (
	"math"
	"testing"

	"github.com/cgi-fr/sigo/pkg/reidentification"
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

	cosine := reidentification.NewCosineSimilarity()

	assert.InDelta(t, 0.7622963735959798, cosine.Compute(X, Y), math.Pow10(-15))
}

func TestEuclideanDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"k1": 0, "k2": 0}
	Y := map[string]float64{"k1": 2, "k2": 0}

	dist := reidentification.NewEuclideanDistance()

	assert.Equal(t, 2.00, dist.Compute(X, Y))
}

func TestManhattanDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 9, "q2": 5}

	dist := reidentification.NewManhattanDistance()

	assert.Equal(t, 10.00, dist.Compute(X, Y))
}

func TestCanberraDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 4, "q2": 4}

	dist := reidentification.NewCanberraDistance()

	assert.Equal(t, float64(2)/float64(3), dist.Compute(X, Y))
}

func TestChebyshevDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 9, "q2": 5}

	dist := reidentification.NewChebyshevDistance()

	assert.Equal(t, 7.00, dist.Compute(X, Y))
}

func TestMinkowskiDistance(t *testing.T) {
	t.Parallel()

	X := map[string]float64{"q1": 2, "q2": 2}
	Y := map[string]float64{"q1": 4, "q2": 4}

	dist := reidentification.NewMinkowskiDistance(6)

	assert.Equal(t, math.Pow(128, 1.00/6.00), dist.Compute(X, Y))
}

package reidentification_test

import (
	"math"
	"testing"

	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()

	X := []float64{14, 6, 18, 52.1, 21, 36.48, 12, 39}
	Y := []float64{14.17, 6, 4, 12.86, 54, 49, 7.2, 27.12}

	cosine := reidentification.CosineSimilarity(X, Y)

	assert.InDelta(t, 0.7622963735959798, cosine, math.Pow10(-15))
}

func TestMapToSlice(t *testing.T) {
	t.Parallel()

	m1 := make(map[string]interface{})
	m1["x"] = 14
	m1["y"] = 6

	m2 := make(map[string]interface{})
	m2["x"] = 1
	m2["y"] = 15

	s1 := reidentification.MapToSlice(m1)
	s2 := reidentification.MapToSlice(m2)

	cosine := reidentification.CosineSimilarity(s1, s2)

	assert.InDelta(t, 0.45418744744022516, cosine, math.Pow10(-15))
}

func TestTopSimilarity(t *testing.T) {
	t.Parallel()

	val1 := []float64{3, 7}
	sim1 := reidentification.NewSimilarity(1, val1, 0.89, "a")
	val2 := []float64{7, 3}
	sim2 := reidentification.NewSimilarity(2, val2, 0.95, "a")
	val3 := []float64{16.67, 18.33}
	sim3 := reidentification.NewSimilarity(3, val3, 0.99, "b")
	val4 := []float64{4.33, 17.67}
	sim4 := reidentification.NewSimilarity(4, val4, 0.80, "c")
	val5 := []float64{16.67, 18.33}
	sim5 := reidentification.NewSimilarity(5, val5, 0.99, "a")

	test := []reidentification.Similarity{sim1, sim2, sim3, sim4, sim5}
	res := reidentification.TopSimilarity(test, 2)
	expected := []reidentification.Similarity{sim3, sim5}
	assert.Equal(t, expected, res)
}

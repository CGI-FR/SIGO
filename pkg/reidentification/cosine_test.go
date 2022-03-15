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

	cosine := reidentification.CosineSimilarity(X, Y)

	assert.InDelta(t, 0.7622963735959798, cosine, math.Pow10(-15))
}

func TestMapInterfaceToFloat(t *testing.T) {
	t.Parallel()

	m1 := make(map[string]interface{})
	m1["x"] = 14
	m1["y"] = 6

	m2 := make(map[string]interface{})
	m2["x"] = 1
	m2["y"] = 15

	s1 := reidentification.MapItoMapF(m1)
	s2 := reidentification.MapItoMapF(m2)

	cosine := reidentification.CosineSimilarity(s1, s2)

	assert.InDelta(t, 0.45418744744022516, cosine, math.Pow10(-15))
}

func TestTopSimilarity(t *testing.T) {
	t.Parallel()

	val1 := map[string]interface{}{"x": 3, "y": 7}
	sim1 := reidentification.NewSimilarity(1, val1, 0.89, "a")
	val2 := map[string]interface{}{"x": 7, "y": 3}
	sim2 := reidentification.NewSimilarity(2, val2, 0.95, "a")
	val3 := map[string]interface{}{"x": 16.67, "y": 18.33}
	sim3 := reidentification.NewSimilarity(3, val3, 0.99, "b")
	val4 := map[string]interface{}{"x": 4.33, "y": 17.67}
	sim4 := reidentification.NewSimilarity(4, val4, 0.80, "c")
	val5 := map[string]interface{}{"x": 16.67, "y": 18.33}
	sim5 := reidentification.NewSimilarity(5, val5, 0.99, "a")

	test := []reidentification.Similarity{sim1, sim2, sim3, sim4, sim5}
	res := reidentification.TopSimilarity(test, 3)
	expected := []reidentification.Similarity{sim3, sim5, sim2}
	assert.Equal(t, expected, res)
}

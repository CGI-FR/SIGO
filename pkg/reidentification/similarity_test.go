package reidentification_test

import (
	"math"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

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

	row := jsonline.NewRow()
	row.Set("x", 4)
	row.Set("y", 4)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})

	x := []float64{3, 7, 16.67, 4.33, 16.67}
	y := []float64{7, 3, 18.33, 17.67, 16}
	z := []string{"a", "a", "b", "c", "a"}
	test := []reidentification.Similarity{}

	for i := range x {
		row1 := jsonline.NewRow()
		row1.Set("x", x[i])
		row1.Set("y", y[i])
		row1.Set("z", z[i])

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i)
		sim.AddSimilarity(record1, []string{"x", "y"}, []string{"z"})
		sim.Compute(record, []string{"x", "y"})

		test = append(test, sim)
	}

	res := reidentification.TopSimilarity(test, 2)

	idE := []int{4, 2}
	xE := []float64{16.67, 16.67}
	yE := []float64{16, 18.33}
	zE := []string{"a", "b"}
	expected := []reidentification.Similarity{}

	for i := range xE {
		rowE := jsonline.NewRow()
		rowE.Set("x", xE[i])
		rowE.Set("y", yE[i])
		rowE.Set("z", zE[i])
		recordE := infra.NewJSONLineRecord(&rowE, &[]string{"x", "y"}, &[]string{"z"})
		simE := reidentification.NewSimilarity(idE[i])
		simE.AddSimilarity(recordE, []string{"x", "y"}, []string{"z"})
		simE.Compute(record, []string{"x", "y"})
		expected = append(expected, simE)
	}

	assert.Equal(t, expected, res)
}

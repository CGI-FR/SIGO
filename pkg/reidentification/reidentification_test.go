package reidentification_test

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestRecover(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("x", 3)
	row.Set("y", 6)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		row1 := jsonline.NewRow()
		row1.Set("x", 3)
		row1.Set("y", 7)
		row1.Set("z", "c")

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i)
		sim.AddSimilarity(record1, []string{"x", "y"}, []string{"z"})
		sim.Compute(record, []string{"x", "y"})

		test = append(test, sim)
	}

	res, risk := reidentification.Recover(test)

	assert.True(t, risk)
	assert.Equal(t, []string{"c"}, res)
}

func TestCountValues(t *testing.T) {
	t.Parallel()

	values := []string{"a", "a", "b", "a", "c", "c", "a", "b"}
	count := reidentification.CountValues(values)

	assert.Equal(t, 4, count["a"])
	assert.Equal(t, 2, count["b"])
	assert.Equal(t, 2, count["c"])
}

func TestRisk(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("x", 11)
	row.Set("y", 9)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		row1 := jsonline.NewRow()
		row1.Set("x", 19.67)
		row1.Set("y", 17.67)
		row1.Set("z", "b")

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i)
		sim.AddSimilarity(record1, []string{"x", "y"}, []string{"z"})
		sim.Compute(record, []string{"x", "y"})

		test = append(test, sim)
	}

	risk := reidentification.Risk(test)

	assert.Equal(t, float64(1), risk)

	test2 := []reidentification.Similarity{}
	z := []string{"a", "b", "b", "b"}

	for i := range z {
		row1 := jsonline.NewRow()
		row1.Set("x", 19.67)
		row1.Set("y", 17.67)
		row1.Set("z", z[i])

		record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})
		sim := reidentification.NewSimilarity(i)
		sim.AddSimilarity(record1, []string{"x", "y"}, []string{"z"})
		sim.Compute(record, []string{"x", "y"})

		test2 = append(test2, sim)
	}

	risk2 := reidentification.Risk(test2)

	assert.Equal(t, float64(0.5), risk2)
}

func TestReidentification(t *testing.T) {
	t.Parallel()

	originalFile, err := os.Open("../../examples/re-identification/data.json")
	assert.Nil(t, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(originalFile), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	var res reidentification.Original

	i := 0

	for original.Next() {
		qi := make(map[string]interface{})

		for _, q := range original.QuasiIdentifer() {
			qi[q] = original.Value().Row()[q]
		}

		var sims reidentification.Similarities

		sigoFile, err := os.Open("../../examples/re-identification/data2-sigo.json")
		assert.Nil(t, err)

		sigo, err := infra.NewJSONLineSource(bufio.NewReader(sigoFile), []string{"x", "y"}, []string{"z"})
		assert.Nil(t, err)

		j := 0

		for sigo.Next() {
			sim := reidentification.NewSimilarity(j)
			sim.AddSimilarity(sigo.Value(), sigo.QuasiIdentifer(), sigo.Sensitive())
			sim.Compute(original.Value(), original.QuasiIdentifer())

			sims.Add(sim)
			j++
		}

		ind := reidentification.NewIndividu(i, qi, sims.Slice())
		i++

		res.Add(ind)
	}

	riskInd := res.Reidenfication(3)
	expected := map[string]interface{}{
		"x": json.Number("20"), "y": json.Number("18"), "sensitive": []string{"b"},
	}
	// assert.Nil(t, riskInd)
	assert.Contains(t, riskInd, expected)

	log.Println(riskInd)
}

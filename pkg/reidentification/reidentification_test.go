package reidentification_test

import (
	"log"
	"strings"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

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

	val := map[string]interface{}{"x": 19.67, "y": 17.67}
	sim1 := reidentification.NewSimilarity(1, val, 0.9999995697145913, "b")
	sim2 := reidentification.NewSimilarity(2, val, 0.9999995697145913, "b")
	sim3 := reidentification.NewSimilarity(3, val, 0.9999995697145913, "b")

	test := []reidentification.Similarity{sim1, sim2, sim3}
	risk := reidentification.Risk(test)

	assert.Equal(t, float64(1), risk)

	val2 := map[string]interface{}{"x": 19.67, "y": 17.67}

	sim11 := reidentification.NewSimilarity(1, val2, 0.99, "a")
	sim22 := reidentification.NewSimilarity(2, val2, 0.99, "b")
	sim33 := reidentification.NewSimilarity(3, val2, 0.99, "b")
	sim44 := reidentification.NewSimilarity(3, val2, 0.99, "b")

	test2 := []reidentification.Similarity{sim11, sim22, sim33, sim44}
	risk2 := reidentification.Risk(test2)

	assert.Equal(t, float64(0.5), risk2)
}

func TestReidentification(t *testing.T) {
	t.Parallel()

	originalData := `{"x": 5, "y": 6, "z":"a"}
					 {"x": 3, "y": 7, "z":"a"}
					 {"x": 4, "y": 4, "z":"c"}
					 {"x": 2, "y": 10, "z":"b"}
					 {"x": 8, "y": 4, "z":"a"}
					 {"x": 8, "y": 10, "z":"a"}
					 {"x": 3, "y": 16, "z":"a"}
					 {"x": 7, "y": 19, "z":"a"}
					 {"x": 6, "y": 18, "z":"a"}
					 {"x": 4, "y": 19, "z":"b"}
					 {"x": 7, "y": 14, "z":"c"}
					 {"x": 10, "y": 14, "z":"c"}
					 {"x": 15, "y": 5, "z":"c"}
					 {"x": 15, "y": 7, "z":"b"}
					 {"x": 11, "y": 9, "z":"b"}
					 {"x": 12, "y": 3, "z":"a"}
					 {"x": 18, "y": 6, "z":"c"}
					 {"x": 14, "y": 6, "z":"c"}
					 {"x": 20, "y": 20, "z":"b"}
					 {"x": 18, "y": 19, "z":"c"}
					 {"x": 20, "y": 18, "z":"b"}
					 {"x": 18, "y": 18, "z":"c"}
					 {"x": 14, "y": 18, "z":"b"}
					 {"x": 19, "y": 15, "z":"b"}`

	sigoData := `{"x":3,"y":7,"z":"b"}
				 {"x":3,"y":7,"z":"a"}
				 {"x":3,"y":7,"z":"c"}
				 {"x":7,"y":6.67,"z":"a"}
				 {"x":7,"y":6.67,"z":"a"}
				 {"x":7,"y":6.67,"z":"a"}
				 {"x":4.33,"y":17.67,"z":"a"}
				 {"x":4.33,"y":17.67,"z":"b"}
				 {"y":17.67,"z":"a","x":4.33}
				 {"x":8,"y":15.67,"z":"c"}
				 {"x":8,"y":15.67,"z":"a"}
				 {"x":8,"y":15.67,"z":"c"}
				 {"x":12.33,"y":6,"z":"b"}
				 {"x":12.33,"y":6,"z":"a"}
				 {"x":12.33,"y":6,"z":"c"}
				 {"y":6,"z":"c","x":16}
				 {"y":6,"z":"b","x":16}
				 {"x":16,"y":6,"z":"c"}
				 {"x":16.67,"y":18.33,"z":"b"}
				 {"y":18.33,"z":"c","x":16.67}
				 {"x":16.67,"y":18.33,"z":"c"}
				 {"y":17.67,"z":"b","x":19.67}
				 {"y":17.67,"z":"b","x":19.67}
				 {"x":19.67,"y":17.67,"z":"b"}`

	original, err := infra.NewJSONLineSource(strings.NewReader(originalData), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	var res []reidentification.Individu
	i := 0

	for original.Next() {
		qi := make(map[string]interface{})

		for _, q := range original.QuasiIdentifer() {
			qi[q] = original.Value().Row()[q]
		}

		individu := reidentification.MapItoMapF(qi)

		var sliceSim []reidentification.Similarity

		sigo, err := infra.NewJSONLineSource(strings.NewReader(sigoData), []string{"x", "y"}, []string{"z"})
		assert.Nil(t, err)

		j := 0

		for sigo.Next() {
			qiSigo := make(map[string]interface{})

			for _, q := range sigo.QuasiIdentifer() {
				qiSigo[q] = sigo.Value().Row()[q]
			}

			individuSigo := reidentification.MapItoMapF(qiSigo)
			score := reidentification.CosineSimilarity(individu, individuSigo)
			sensitive := sigo.Value().Row()[sigo.Sensitive()[0]]
			sim := reidentification.NewSimilarity(j, qiSigo, score, sensitive.(string))

			sliceSim = append(sliceSim, sim)
			j++
		}

		ind := reidentification.NewIndividu(i, qi, sliceSim)
		i++

		res = append(res, ind)
	}

	log.Println("values", res[20].Values())
	top := reidentification.TopSimilarity(res[20].Similarities(), 3)
	log.Println("top", top)
	log.Println(reidentification.Risk(top))
	log.Println(res[20].Reidenfication(3))

	// assert.Equal(t, 19.00, q.Q1)
}

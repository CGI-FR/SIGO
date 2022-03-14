package reidentification_test

import (
	"log"
	"strings"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

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

	for original.Next() {
		data := make(map[string]interface{})

		for _, q := range original.QuasiIdentifer() {
			data[q] = original.Value().Row()[q]
		}

		individu := reidentification.MapToSlice(data)

		var sliceSim []reidentification.Similarity

		sigo, err := infra.NewJSONLineSource(strings.NewReader(sigoData), []string{"x", "y"}, []string{"z"})
		assert.Nil(t, err)

		for sigo.Next() {
			i := 0
			dataSigo := make(map[string]interface{})

			for _, q := range sigo.QuasiIdentifer() {
				dataSigo[q] = sigo.Value().Row()[q]
			}

			individuSigo := reidentification.MapToSlice(dataSigo)
			score := reidentification.CosineSimilarity(individu, individuSigo)
			sensitive := sigo.Value().Row()[sigo.Sensitive()[0]]
			sim := reidentification.NewSimilarity(i, individuSigo, score, sensitive.(string))

			sliceSim = append(sliceSim, sim)
			i++
		}

		ind := reidentification.NewIndividu(individu, sliceSim)
		res = append(res, ind)
	}

	log.Println(res[20].Values())
	log.Println(reidentification.TopSimilarity(res[20].Similarities(), 5))

	// assert.Equal(t, 19.00, q.Q1)
}

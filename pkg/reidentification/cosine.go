package reidentification

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
)

func CosineSimilarity(x, y []float64) float64 {
	var dotProduct, X, Y float64

	//nolint: gomnd
	for i := range x {
		dotProduct += x[i] * y[i]
		X += math.Pow(x[i], 2)
		Y += math.Pow(y[i], 2)
	}

	return dotProduct / (math.Sqrt(X) * math.Sqrt(Y))
}

func MapToSlice(m map[string]interface{}) (values []float64) {
	for _, value := range m {
		var val float64
		switch t := value.(type) {
		case int:
			val = float64(t)
		case string:
			//nolint: gomnd
			val, _ = strconv.ParseFloat(t, 64)
		case float32:
			val = float64(t)
		case json.Number:
			val, _ = t.Float64()
		}

		values = append(values, val)
	}

	return values
}

func TopSimilarity(sims []Similarity, n int) (res []Similarity) {
	type tmp struct {
		individu  []float64
		score     float64
		sensitive string
	}

	m := make(map[float64][]int)
	scores := []float64{}
	mapTmp := make(map[int]tmp)

	for _, sim := range sims {
		var t tmp

		scores = append(scores, sim.score)
		m[sim.score] = append(m[sim.score], sim.id)
		t.individu = sim.individu
		t.score = sim.score
		t.sensitive = sim.sensitive
		mapTmp[sim.id] = t
	}

	scores = removeDuplicate(scores)
	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	count := 0

	for _, k := range scores {
		for _, i := range m[k] {
			ind := Similarity{id: i, individu: mapTmp[i].individu, score: k, sensitive: mapTmp[i].sensitive}
			count++

			res = append(res, ind)
		}

		if count == n {
			break
		}
	}

	return res
}

func removeDuplicate(floatSlice []float64) []float64 {
	keys := make(map[float64]bool)
	list := []float64{}

	for _, val := range floatSlice {
		if _, value := keys[val]; !value {
			keys[val] = true

			list = append(list, val)
		}
	}

	return list
}

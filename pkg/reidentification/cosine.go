package reidentification

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
)

func CosineSimilarity(x, y map[string]float64) float64 {
	var dotProduct, X, Y float64

	//nolint: gomnd
	for key := range x {
		dotProduct += x[key] * y[key]
		X += math.Pow(x[key], 2)
		Y += math.Pow(y[key], 2)
	}

	return dotProduct / (math.Sqrt(X) * math.Sqrt(Y))
}

func MapItoMapF(m map[string]interface{}) map[string]float64 {
	mFloat := make(map[string]float64)

	for key, value := range m {
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

		mFloat[key] = val
	}

	return mFloat
}

func TopSimilarity(sims []Similarity, n int) (res []Similarity) {
	type tmp struct {
		individu  map[string]interface{}
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
		t.individu = sim.qi
		t.score = sim.score
		t.sensitive = sim.sensitive
		mapTmp[sim.id] = t
	}

	scores = removeDuplicate(scores)
	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	count := 0

	for _, k := range scores {
		for _, i := range m[k] {
			ind := Similarity{id: i, qi: mapTmp[i].individu, score: k, sensitive: mapTmp[i].sensitive}
			count++

			res = append(res, ind)

			if count == n {
				break
			}
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

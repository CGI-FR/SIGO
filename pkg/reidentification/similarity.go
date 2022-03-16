package reidentification

import (
	"encoding/json"
	"reflect"
	"sort"
	"strconv"
)

type Metric interface {
	Compute(map[string]float64, map[string]float64) float64
}

type Similarities struct {
	slice  []Similarity
	metric Metric
}

func NewSimilarities(m Metric) Similarities {
	return Similarities{slice: []Similarity{}, metric: m}
}

func (s *Similarities) Add(sim Similarity) {
	s.slice = append(s.slice, sim)
}

func (s Similarities) Slice() []Similarity {
	return s.slice
}

func (s Similarities) Metric() Metric {
	return s.metric
}

func (s Similarities) TopSimilarity(n int) (res Similarities) {
	type tmp struct {
		individu  map[string]interface{}
		score     float64
		sensitive []string
	}

	m := make(map[float64][]int)
	scores := []float64{}
	mapTmp := make(map[int]tmp)

	for _, sim := range s.slice {
		var t tmp

		scores = append(scores, sim.score)
		m[sim.score] = append(m[sim.score], sim.id)
		t.individu = sim.qi
		t.score = sim.score
		t.sensitive = sim.sensitive
		mapTmp[sim.id] = t
	}

	scores = removeDuplicate(scores)

	switch reflect.TypeOf(s.metric).String() {
	case "reidentification.Cosine":
		sort.Sort(sort.Reverse(sort.Float64Slice(scores)))
	default:
		sort.Sort(sort.Float64Slice(scores))
	}

	count := 0

	for _, k := range scores {
		for _, i := range m[k] {
			ind := Similarity{id: i, qi: mapTmp[i].individu, score: k, sensitive: mapTmp[i].sensitive}
			count++

			res.slice = append(res.slice, ind)

			if count == n {
				break
			}
		}

		if count == n {
			break
		}
	}

	res.metric = s.metric

	return res
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
		case float64:
			val = t
		}

		mFloat[key] = val
	}

	return mFloat
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

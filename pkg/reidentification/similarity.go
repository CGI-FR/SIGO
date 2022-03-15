package reidentification

import (
	"encoding/json"
	"sort"
	"strconv"

	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Similarities struct {
	sims []Similarity
}

func NewSimilarities() Similarities {
	return Similarities{sims: []Similarity{}}
}

func (s *Similarities) Add(sim Similarity) {
	s.sims = append(s.sims, sim)
}

func (s Similarities) Slice() []Similarity {
	return s.sims
}

type Similarity struct {
	id        int
	qi        map[string]interface{}
	score     float64
	sensitive []string
}

func NewSimilarity(id int) Similarity {
	return Similarity{id: id, qi: make(map[string]interface{}), score: 0, sensitive: []string{}}
}

func (sim *Similarity) AddSimilarity(ind sigo.Record, qid []string, s []string) {
	for _, q := range qid {
		sim.qi[q] = ind.Row()[q]
	}

	for i := range s {
		sim.sensitive = append(sim.sensitive, ind.Row()[s[i]].(string))
	}
}

func (sim *Similarity) Compute(ind sigo.Record, qid []string) {
	x := make(map[string]interface{})

	for _, q := range qid {
		x[q] = ind.Row()[q]
	}

	X := MapItoMapF(x)
	Y := MapItoMapF(sim.qi)

	sim.score = CosineSimilarity(X, Y)
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

func TopSimilarity(sims []Similarity, n int) (res []Similarity) {
	type tmp struct {
		individu  map[string]interface{}
		score     float64
		sensitive []string
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

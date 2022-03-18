// Copyright (C) 2022 CGI France
//
// This file is part of SIGO.
//
// SIGO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SIGO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SIGO.  If not, see <http://www.gnu.org/licenses/>.

package reidentification

import (
	"reflect"
	"sort"
)

type Similarity struct {
	id        int
	qi        map[string]interface{}
	score     float64
	sensitive []string
}

func NewSimilarity(id int, ind map[string]interface{}, qid []string, s []string) Similarity {
	list := []string{}

	for _, i := range s {
		list = append(list, ind[i].(string))
	}

	return Similarity{id: id, qi: ind, score: 0, sensitive: list}
}

func (sim *Similarity) AddScore(score float64) {
	sim.score = score
}

type Similarities struct {
	slice  []Similarity
	metric Distance
}

func NewSimilarities(m Distance) Similarities {
	return Similarities{slice: []Similarity{}, metric: m}
}

func (s *Similarities) Add(sim Similarity) {
	s.slice = append(s.slice, sim)
}

func (s Similarities) Slice() []Similarity {
	return s.slice
}

func (s Similarities) Metric() Distance {
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

	scores = RemoveDuplicate(scores)

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

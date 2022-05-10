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
	"sort"
)

type Similarity struct {
	id        int
	qi        map[string]interface{}
	score     float64
	sensitive []string
}

// NewSimilarities instantiates a Similarity object, which associates a score with anonymized data.
// This score represents the distance between anonymized data and original data.
func NewSimilarity(id int, ind map[string]interface{}, qid []string, s []string) Similarity {
	list := []string{}

	for _, i := range s {
		list = append(list, ind[i].(string))
	}

	return Similarity{id: id, qi: ind, score: 0, sensitive: list}
}

// AddScore adds score to similarity.
func (sim *Similarity) AddScore(score float64) {
	sim.score = score
}

type Similarities struct {
	slice  []Similarity
	metric string
}

// NewSimilarities instantiates a Similarities object.
func NewSimilarities(m string) Similarities {
	return Similarities{slice: []Similarity{}, metric: m}
}

// Add adds a similarity to Similarities slice.
func (s *Similarities) Add(sim Similarity) {
	s.slice = append(s.slice, sim)
}

// Slice returns the Similarities slice.
func (s Similarities) Slice() []Similarity {
	return s.slice
}

// Metric returns the metric use to calculate the similarities.
func (s Similarities) Metric() string {
	return s.metric
}

// TopSimilarity returns the best score.
func (s Similarities) TopSimilarity() (res Similarity) {
	type tmp struct {
		individu  map[string]interface{}
		score     float64
		sensitive []string
	}

	// map containing for each distance score the identifier of the individual associated to this score
	m := make(map[float64]int)
	scores := []float64{}
	mapTmp := make(map[int]tmp)

	for _, sim := range s.slice {
		var t tmp

		scores = append(scores, sim.score)
		m[sim.score] = sim.id
		t.individu = sim.qi
		t.score = sim.score
		t.sensitive = sim.sensitive
		mapTmp[sim.id] = t
	}

	switch s.metric {
	case "cosine":
		// Score closest to 1 for the cosine distance
		sort.Sort(sort.Reverse(sort.Float64Slice(scores)))
	default:
		// Scores closest to 0 for the other distance
		sort.Sort(sort.Float64Slice(scores))
	}

	// best score
	top := scores[0]

	// recover the individual associated with the best score
	i := m[top]

	res.id = i
	res.qi = mapTmp[i].individu
	res.score = top
	res.sensitive = mapTmp[i].sensitive

	return res
}

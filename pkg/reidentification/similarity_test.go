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

package reidentification_test

import (
	"testing"

	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestTopSimilarity(t *testing.T) {
	t.Parallel()

	x := []float64{3, 7, 16.67, 4.33, 16.67}
	y := []float64{7, 3, 18.33, 17.67, 16}
	z := []string{"a", "a", "b", "c", "a"}
	scores := []float64{0.8, 0.5, 0.9, 0.6, 0.9}
	test := reidentification.NewSimilarities(reidentification.NewCosineSimilarity())

	for i := range x {
		record := make(map[string]interface{})
		record["x"] = x[i]
		record["y"] = y[i]
		record["z"] = z[i]

		sim := reidentification.NewSimilarity(i, record, []string{"x", "y"}, []string{"z"})
		sim.AddScore(scores[i])

		test.Add(sim)
	}

	res := test.TopSimilarity(2)

	idE := []int{2, 4}
	xE := []float64{16.67, 16.67}
	yE := []float64{18.33, 16}
	zE := []string{"b", "a"}
	scoresE := []float64{0.9, 0.9}
	expected := reidentification.NewSimilarities(reidentification.NewCosineSimilarity())

	for i := range xE {
		recordE := make(map[string]interface{})
		recordE["x"] = xE[i]
		recordE["y"] = yE[i]
		recordE["z"] = zE[i]

		simE := reidentification.NewSimilarity(idE[i], recordE, []string{"x", "y"}, []string{"z"})
		simE.AddScore(scoresE[i])
		expected.Add(simE)
	}

	assert.Equal(t, expected, res)
}

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
	"math"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestMapInterfaceToFloat(t *testing.T) {
	t.Parallel()

	m1 := make(map[string]interface{})
	m1["x"] = 14
	m1["y"] = 6

	m2 := make(map[string]interface{})
	m2["x"] = 1
	m2["y"] = 15

	s1 := reidentification.MapItoMapF(m1)
	s2 := reidentification.MapItoMapF(m2)

	dist := reidentification.Cosine(s1, s2)

	assert.InDelta(t, 0.45418744744022516, dist, math.Pow10(-15))
}

func TestRecover(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("x", 3)
	row.Set("y", 6)
	// record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		record1 := make(map[string]interface{})
		record1["x"] = 3
		record1["y"] = 7
		record1["z"] = "c"

		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		// sim.ComputeSimilarity(record, []string{"x", "y"}, reidentification.NewCosineSimilarity())

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
	// record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{})
	test := []reidentification.Similarity{}

	for i := 0; i < 3; i++ {
		record1 := make(map[string]interface{})
		record1["x"] = 19.67
		record1["y"] = 17.67
		record1["z"] = "b"

		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		test = append(test, sim)
	}

	risk := reidentification.Risk(test)

	assert.Equal(t, float64(1), risk)

	test2 := []reidentification.Similarity{}
	z := []string{"a", "b", "b", "b"}

	for i := range z {
		record1 := make(map[string]interface{})
		record1["x"] = 19.67
		record1["y"] = 17.67
		record1["z"] = z[i]

		sim := reidentification.NewSimilarity(i, record1, []string{"x", "y"}, []string{"z"})
		test2 = append(test2, sim)
	}

	risk2 := reidentification.Risk(test2)

	assert.Equal(t, float64(0.5), risk2)
}

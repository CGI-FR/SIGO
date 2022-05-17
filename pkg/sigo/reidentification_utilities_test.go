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

package sigo_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/cgi-fr/sigo/pkg/sigo"
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

	s1 := sigo.MapItoMapF(m1)
	s2 := sigo.MapItoMapF(m2)

	dist := sigo.Cosine(s1, s2)

	assert.InDelta(t, 0.45418744744022516, dist, math.Pow10(-15))
}

func TestCountValues(t *testing.T) {
	t.Parallel()

	values := []string{"a", "a", "b", "a", "c", "c", "a", "b"}
	count := sigo.CountValues(values)

	assert.Equal(t, 4, count["a"])
	assert.Equal(t, 2, count["b"])
	assert.Equal(t, 2, count["c"])
}

func TestListValues(t *testing.T) {
	t.Parallel()

	data := make([]map[string]interface{}, 3)

	vals0 := make(map[string]interface{})
	vals0["x"] = 2.12
	vals0["y"] = 4.5
	vals0["z"] = "a"

	vals1 := make(map[string]interface{})
	vals1["x"] = 4.36
	vals1["y"] = 8.75
	vals1["z"] = "b"

	vals2 := make(map[string]interface{})
	vals2["x"] = 12.17
	vals2["y"] = 3.96
	vals2["z"] = "c"

	data[0] = vals0
	data[1] = vals1
	data[2] = vals2

	res := sigo.ListValues(data, []string{"z"})

	expected := make(map[string][]interface{})
	expected["x"] = []interface{}{2.12, 4.36, 12.17}
	expected["y"] = []interface{}{4.5, 8.75, 3.96}

	assert.Equal(t, expected, res)
}

func TestScaledData(t *testing.T) {
	t.Parallel()

	data := make([]map[string]interface{}, 3)
	x := []float64{10, 30, 50}
	y := []float64{20, 40, 50}
	z := []string{"a", "b", "c"}

	for i := range x {
		vals := make(map[string]interface{})
		vals["x"] = x[i]
		vals["y"] = y[i]
		vals["z"] = z[i]

		data[i] = vals
	}

	dataScaled := make([]map[string]interface{}, 3)
	xScaled := []float64{-1.224744871391589, 0, 1.224744871391589}
	yScaled := []float64{-1.3365734230692703, 0.26699397113501305, 1.0687776682371546}
	zScaled := []string{"a", "b", "c"}

	for i := range xScaled {
		vals := make(map[string]interface{})
		vals["x"] = xScaled[i]
		vals["y"] = yScaled[i]
		vals["z"] = zScaled[i]

		dataScaled[i] = vals
	}

	res := sigo.ScaleData(data, []string{"z"})

	assert.Equal(t, dataScaled, res)
}

func TestRoundFloat(t *testing.T) {
	t.Parallel()

	v := 0.957895
	v2 := 0.6058987922965455

	res := sigo.RoundFloat(v, 4)
	res2 := sigo.RoundFloat(v2, 4)

	assert.Equal(t, 0.9579, res)
	assert.Equal(t, 0.6059, res2)
}

func TestUnique(t *testing.T) {
	t.Parallel()

	slice := make([]map[string]interface{}, 3)

	data := make(map[string]interface{})
	data["original"] = 0
	data["x"] = json.Number("7")
	data["y"] = json.Number("6.67")
	data["z"] = "a"

	slice[0] = data
	slice[1] = data
	slice[2] = data

	res := sigo.Unique(slice, []string{"x", "y"})

	assert.True(t, res)
}

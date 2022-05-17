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

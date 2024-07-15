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
	"testing"

	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestCountUniqueValues(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	source := sigo.NewAnalyzer(qi)

	source.Add(createRow(4, 1, qi))
	source.Add(createRow(3, 2, qi))
	source.Add(createRow(4, 3, qi))

	res := source.CountUniqueValues()

	assert.Equal(t, 2, res["x"])
	assert.Equal(t, 3, res["y"])
}

func TestOrderMap(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	source := sigo.NewAnalyzer(qi)

	source.Add(createRow(4, 1, qi))
	source.Add(createRow(3, 2, qi))
	source.Add(createRow(4, 3, qi))

	res := source.Order()

	assert.Equal(t, "y", res[0])
	assert.Equal(t, "x", res[1])
}

func TestDimension(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	source := sigo.NewAnalyzer(qi)

	source.Add(createRow(4, 1, qi))
	source.Add(createRow(3, 2, qi))
	source.Add(createRow(4, 3, qi))

	res := source.Dimension(0)

	assert.Equal(t, 1, res)
}

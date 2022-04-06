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
	"strings"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestAddRecordToAnalyzer(t *testing.T) {
	t.Parallel()

	sourceText := `{"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	assert.Nil(t, err)

	analyzer := sigo.New(source)

	row := jsonline.NewRow()
	row.Set("x", 3)
	row.Set("y", 2)
	row.Set("z", "bar")
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{"foo"})

	analyzer.Add(record)

	assert.Equal(t, analyzer.Values("x"), []float64{3})
	assert.Equal(t, analyzer.Values("y"), []float64{2})
}

func TestCountUniqueValues(t *testing.T) {
	t.Parallel()

	sourceText := `{"x":4, "y":1}
				   {"x":3, "y":2}
				   {"x":4, "y":3}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{})
	assert.Nil(t, err)

	analyzer := sigo.New(source)

	qi := []string{"x", "y"}
	analyzer.Add(createRow(4, 1, qi))
	analyzer.Add(createRow(3, 2, qi))
	analyzer.Add(createRow(4, 3, qi))

	res := analyzer.CountUniqueValues()

	assert.Equal(t, 2, res["x"])
	assert.Equal(t, 3, res["y"])
}

func TestOrderQI(t *testing.T) {
	t.Parallel()

	sourceText := `{"x":1, "y":1}
				   {"x":2, "y":2}
				   {"x":1, "y":3}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{})
	assert.Nil(t, err)

	analyzer := sigo.New(source)

	qi := []string{"x", "y"}
	analyzer.Add(createRow(1, 1, qi))
	analyzer.Add(createRow(2, 2, qi))
	analyzer.Add(createRow(1, 3, qi))

	unique := analyzer.CountUniqueValues()
	res := sigo.Order(unique)

	assert.Equal(t, []string{"y", "x"}, res)
}

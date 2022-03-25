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

func TestAggregationAnonymizer(t *testing.T) {
	t.Parallel()

	row1 := jsonline.NewRow()
	row1.Set("x", 0)
	row1.Set("y", 9)
	row1.Set("z", "a")
	record1 := infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{"z"})

	tree := sigo.NewKDTree(2, 3, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	node.Add(record1)

	row2 := jsonline.NewRow()
	row2.Set("x", 1)
	row2.Set("y", 3)
	row2.Set("z", "b")
	record2 := infra.NewJSONLineRecord(&row2, &[]string{"x", "y"}, &[]string{"z"})
	node.Add(record2)

	row3 := jsonline.NewRow()
	row3.Set("x", 4)
	row3.Set("y", 8)
	row3.Set("z", "c")
	record3 := infra.NewJSONLineRecord(&row3, &[]string{"x", "y"}, &[]string{"z"})
	node.Add(record3)

	anonymizer := sigo.NewAggregationAnonymizer("mean")
	anonymizedRecord := anonymizer.Anonymize(record1, node.Clusters()[0], []string{"x", "y"}, []string{"z"})
	expected := map[string]interface{}{"x": 1.67, "y": 6.67, "z": "a"}

	assert.Equal(t, expected, anonymizedRecord.Row())

	anonymizer = sigo.NewAggregationAnonymizer("median")
	anonymizedRecord = anonymizer.Anonymize(record1, node.Clusters()[0], []string{"x", "y"}, []string{"z"})
	expected = map[string]interface{}{"x": 1.00, "y": 8.00, "z": "a"}

	assert.Equal(t, expected, anonymizedRecord.Row())
}

//nolint: funlen
func TestTopBottomCodingAnonymizer(t *testing.T) {
	t.Parallel()

	tree := sigo.NewKDTree(7, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)

	row1 := jsonline.NewRow()
	row1.Set("x", 3)
	row1.Set("y", 5)
	node.Add(infra.NewJSONLineRecord(&row1, &[]string{"x", "y"}, &[]string{}))

	row2 := jsonline.NewRow()
	row2.Set("x", 5)
	row2.Set("y", 3)
	record2 := infra.NewJSONLineRecord(&row2, &[]string{"x", "y"}, &[]string{})
	node.Add(record2)

	row3 := jsonline.NewRow()
	row3.Set("x", 6)
	row3.Set("y", 5)
	node.Add(infra.NewJSONLineRecord(&row3, &[]string{"x", "y"}, &[]string{}))

	row4 := jsonline.NewRow()
	row4.Set("x", 9)
	row4.Set("y", 10)
	node.Add(infra.NewJSONLineRecord(&row4, &[]string{"x", "y"}, &[]string{}))

	row5 := jsonline.NewRow()
	row5.Set("x", 10)
	row5.Set("y", 30)
	node.Add(infra.NewJSONLineRecord(&row5, &[]string{"x", "y"}, &[]string{}))

	row6 := jsonline.NewRow()
	row6.Set("x", 11)
	row6.Set("y", 11)
	node.Add(infra.NewJSONLineRecord(&row6, &[]string{"x", "y"}, &[]string{}))

	row7 := jsonline.NewRow()
	row7.Set("x", 12)
	row7.Set("y", 11)
	node.Add(infra.NewJSONLineRecord(&row7, &[]string{"x", "y"}, &[]string{}))

	row8 := jsonline.NewRow()
	row8.Set("x", 48)
	row8.Set("y", 12)
	record8 := infra.NewJSONLineRecord(&row8, &[]string{"x", "y"}, &[]string{})
	node.Add(record8)

	anonymizer := sigo.NewCodingAnonymizer()

	anonymizedRecord := anonymizer.Anonymize(record2, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected := map[string]interface{}{"x": 5.50, "y": 5.00}

	assert.Equal(t, expected, anonymizedRecord.Row())

	anonymizedRecord = anonymizer.Anonymize(record8, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected = map[string]interface{}{"x": 11.50, "y": 11.50}

	assert.Equal(t, expected, anonymizedRecord.Row())
}

func TestRandomNoiseAnonymizer(t *testing.T) {
	t.Parallel()

	sourceText := `{"x": 0, "y": 0}
				   {"x": 0, "y": 1}
				   {"x": 0, "y": 12}
				   {"x": 1, "y": 1}
				   {"x": 1, "y": 2}
				   {"x": 1, "y": 20}
				   {"x": 2, "y": 1}
				   {"x": 3, "y": 5}
				   {"x": 5, "y": 3}
				   {"x": 6, "y": 5}
				   {"x": 9, "y": 10}
				   {"x": 10, "y": 30}
				   {"x": 11, "y": 11}
				   {"x": 12, "y": 11}
				   {"x": 48, "y": 12}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{})
	assert.Nil(t, err)

	result := []map[string]interface{}{}
	sink := infra.NewSliceDictionariesSink(&result)
	err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), 3, 1, 2, sigo.NewNoiseAnonymizer("gaussian"),
		sink, sigo.NewNoDebugger())
	assert.Nil(t, err)

	for i := 0; i < 4; i++ {
		assert.GreaterOrEqual(t, result[i]["x"], 0.00)
		assert.LessOrEqual(t, result[i]["x"], 2.00)
		assert.GreaterOrEqual(t, result[i]["y"], 0.00)
		assert.LessOrEqual(t, result[i]["y"], 1.00)
	}

	for i := 4; i < 7; i++ {
		assert.GreaterOrEqual(t, result[i]["x"], 0.00)
		assert.LessOrEqual(t, result[i]["x"], 1.00)
		assert.GreaterOrEqual(t, result[i]["y"], 2.00)
		assert.LessOrEqual(t, result[i]["y"], 20.00)
	}

	for i := 7; i < 11; i++ {
		assert.GreaterOrEqual(t, result[i]["x"], 3.00)
		assert.LessOrEqual(t, result[i]["x"], 9.00)
		assert.GreaterOrEqual(t, result[i]["y"], 3.00)
		assert.LessOrEqual(t, result[i]["y"], 10.00)
	}

	for i := 11; i < 15; i++ {
		assert.GreaterOrEqual(t, result[i]["x"], 10.00)
		assert.LessOrEqual(t, result[i]["x"], 48.00)
		assert.GreaterOrEqual(t, result[i]["y"], 11.00)
		assert.LessOrEqual(t, result[i]["y"], 30.00)
	}
}

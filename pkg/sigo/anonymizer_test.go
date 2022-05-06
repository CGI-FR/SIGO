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
	"log"
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

func BenchmarkAggregationAnonymizer(b *testing.B) {
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

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		anonymizer.Anonymize(record1, node.Clusters()[0], []string{"x", "y"}, []string{"z"})
	}
}

func TestTopBottomCodingAnonymizer(t *testing.T) {
	t.Parallel()

	tree := sigo.NewKDTree(7, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	node.Add(createRow(3, 5, qi))
	record2 := createRow(5, 3, qi)
	node.Add(record2)
	node.Add(createRow(6, 5, qi))
	node.Add(createRow(9, 10, qi))
	node.Add(createRow(10, 30, qi))
	node.Add(createRow(11, 11, qi))
	node.Add(createRow(12, 11, qi))
	record8 := createRow(48, 12, qi)
	node.Add(record8)

	anonymizer := sigo.NewCodingAnonymizer()

	anonymizedRecord := anonymizer.Anonymize(record2, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected := map[string]interface{}{"x": 5.50, "y": 5.00}

	assert.Equal(t, expected, anonymizedRecord.Row())

	anonymizedRecord = anonymizer.Anonymize(record8, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected = map[string]interface{}{"x": 11.50, "y": 11.50}

	assert.Equal(t, expected, anonymizedRecord.Row())
}

func BenchmarkTopBottomCodingAnonymizer(b *testing.B) {
	tree := sigo.NewKDTree(7, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	node.Add(createRow(3, 5, qi))
	record2 := createRow(5, 3, qi)
	node.Add(record2)
	node.Add(createRow(6, 5, qi))
	node.Add(createRow(9, 10, qi))
	node.Add(createRow(10, 30, qi))
	node.Add(createRow(11, 11, qi))
	node.Add(createRow(12, 11, qi))
	record8 := createRow(48, 12, qi)
	node.Add(record8)

	anonymizer := sigo.NewCodingAnonymizer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		anonymizer.Anonymize(record2, node.Clusters()[0], []string{"x", "y"}, []string{})
	}
}

func TestRandomNoiseAnonymizer(t *testing.T) {
	t.Parallel()

	tree := sigo.NewKDTree(3, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	record := createRow(3, 5, qi)
	node.Add(record)
	node.Add(createRow(5, 3, qi))
	node.Add(createRow(6, 5, qi))
	node.Add(createRow(9, 10, qi))

	anonymizer := sigo.NewNoiseAnonymizer("gaussian")

	anonymizedRecord := anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})

	assert.GreaterOrEqual(t, anonymizedRecord.Row()["x"], 3.00)
	assert.LessOrEqual(t, anonymizedRecord.Row()["x"], 9.00)
	assert.GreaterOrEqual(t, anonymizedRecord.Row()["y"], 3.00)
	assert.LessOrEqual(t, anonymizedRecord.Row()["y"], 10.00)
}

func BenchmarkRandomNoiseAnonymizer(b *testing.B) {
	tree := sigo.NewKDTree(3, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	record := createRow(3, 5, qi)
	node.Add(record)
	node.Add(createRow(5, 3, qi))
	node.Add(createRow(6, 5, qi))
	node.Add(createRow(9, 10, qi))

	anonymizer := sigo.NewNoiseAnonymizer("gaussian")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})
	}
}

func TestSwapAnonymizer(t *testing.T) {
	t.Parallel()

	tree := sigo.NewKDTree(3, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	record := createRow(1, 4, qi)
	node.Add(record)
	node.Add(createRow(2, 5, qi))
	node.Add(createRow(3, 6, qi))

	anonymizer := sigo.NewSwapAnonymizer()

	anonymizedRecord := anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})

	log.Println(anonymizedRecord.Row()["x"])
	log.Println(anonymizedRecord.Row()["y"])

	assert.Contains(t, []float64{1, 2, 3}, anonymizedRecord.Row()["x"])
	assert.Contains(t, []float64{4, 5, 6}, anonymizedRecord.Row()["y"])
}

func BenchmarkSwapAnonymizer(b *testing.B) {
	tree := sigo.NewKDTree(3, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	record := createRow(1, 4, qi)
	node.Add(record)
	node.Add(createRow(2, 5, qi))
	node.Add(createRow(3, 6, qi))

	anonymizer := sigo.NewSwapAnonymizer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})
	}
}

func createRow(x, y float64, qi []string) infra.JSONLineRecord {
	row := jsonline.NewRow()
	row.Set("x", x)
	row.Set("y", y)

	return infra.NewJSONLineRecord(&row, &qi, &[]string{})
}

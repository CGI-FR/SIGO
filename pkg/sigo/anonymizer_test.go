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
	"log"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestAggregationAnonymizer(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	s := []string{"z"}

	tree := sigo.NewKDTree(2, 3, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)

	record1 := createRow(0, 9, qi, "a", s)
	node.Add(record1)
	node.Add(createRow(1, 3, qi, "b", s))
	node.Add(createRow(4, 8, qi, "c", s))

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
	qi := []string{"x", "y"}
	s := []string{"z"}

	tree := sigo.NewKDTree(2, 3, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)

	record1 := createRow(0, 9, qi, "a", s)
	node.Add(record1)
	node.Add(createRow(1, 3, qi, "b", s))
	node.Add(createRow(4, 8, qi, "c", s))

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

	node.Add(createRow(3, 5, qi, "", []string{}))
	record2 := createRow(5, 3, qi, "", []string{})
	node.Add(record2)
	node.Add(createRow(6, 5, qi, "", []string{}))
	node.Add(createRow(9, 10, qi, "", []string{}))
	node.Add(createRow(10, 30, qi, "", []string{}))
	node.Add(createRow(11, 11, qi, "", []string{}))
	node.Add(createRow(12, 11, qi, "", []string{}))
	record8 := createRow(48, 12, qi, "", []string{})
	node.Add(record8)

	anonymizer := sigo.NewCodingAnonymizer()

	anonymizedRecord := anonymizer.Anonymize(record2, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected := map[string]interface{}{"x": 5.50, "y": 5.00, "z": ""}

	assert.Equal(t, expected, anonymizedRecord.Row())

	anonymizedRecord = anonymizer.Anonymize(record8, node.Clusters()[0], []string{"x", "y"}, []string{})
	expected = map[string]interface{}{"x": 11.50, "y": 11.50, "z": ""}

	assert.Equal(t, expected, anonymizedRecord.Row())
}

func BenchmarkTopBottomCodingAnonymizer(b *testing.B) {
	tree := sigo.NewKDTree(7, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	node.Add(createRow(3, 5, qi, "", []string{}))
	record2 := createRow(5, 3, qi, "", []string{})
	node.Add(record2)
	node.Add(createRow(6, 5, qi, "", []string{}))
	node.Add(createRow(9, 10, qi, "", []string{}))
	node.Add(createRow(10, 30, qi, "", []string{}))
	node.Add(createRow(11, 11, qi, "", []string{}))
	node.Add(createRow(12, 11, qi, "", []string{}))
	record8 := createRow(48, 12, qi, "", []string{})
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

	record := createRow(3, 5, qi, "", []string{})
	node.Add(record)
	node.Add(createRow(5, 3, qi, "", []string{}))
	node.Add(createRow(6, 5, qi, "", []string{}))
	node.Add(createRow(9, 10, qi, "", []string{}))

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

	record := createRow(3, 5, qi, "", []string{})
	node.Add(record)
	node.Add(createRow(5, 3, qi, "", []string{}))
	node.Add(createRow(6, 5, qi, "", []string{}))
	node.Add(createRow(9, 10, qi, "", []string{}))

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

	record := createRow(1, 4, qi, "", []string{})
	node.Add(record)
	node.Add(createRow(2, 5, qi, "", []string{}))
	node.Add(createRow(3, 6, qi, "", []string{}))

	anonymizer := sigo.NewSwapAnonymizer()

	anonymizedRecord := anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})

	assert.Contains(t, []float64{1, 2, 3}, anonymizedRecord.Row()["x"])
	assert.Contains(t, []float64{4, 5, 6}, anonymizedRecord.Row()["y"])
}

func BenchmarkSwapAnonymizer(b *testing.B) {
	tree := sigo.NewKDTree(3, 1, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)
	qi := []string{"x", "y"}

	record := createRow(1, 4, qi, "", []string{})
	node.Add(record)
	node.Add(createRow(2, 5, qi, "", []string{}))
	node.Add(createRow(3, 6, qi, "", []string{}))

	anonymizer := sigo.NewSwapAnonymizer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		anonymizer.Anonymize(record, node.Clusters()[0], []string{"x", "y"}, []string{})
	}
}

func TestReidentification(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	s := []string{"original", "z"}

	tree := sigo.NewKDTree(3, 2, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)

	record1 := createRowReidentification(5, 6, qi, json.Number("1"), s, "")
	node.Add(record1)

	record2 := createRowReidentification(7, 6.67, qi, json.Number("0"), s, "a")
	node.Add(record2)
	node.Add(createRowReidentification(7, 6.67, qi, json.Number("0"), s, "b"))
	node.Add(createRowReidentification(8, 4, qi, json.Number("1"), s, ""))
	node.Add(createRowReidentification(7, 6.67, qi, json.Number("0"), s, "c"))
	node.Add(createRowReidentification(8, 10, qi, json.Number("1"), s, ""))

	anonymizer := sigo.NewReidentification()
	anonymizedRecord := anonymizer.Anonymize(record1, node.Clusters()[0], []string{"x", "y"}, []string{"original", "z"})
	expected := map[string]interface{}{
		"original": json.Number("1"),
		"x":        5.00,
		"y":        6.00,
		"z":        "",
	}

	assert.Equal(t, expected, anonymizedRecord.Row())
}

func TestReidentification2(t *testing.T) {
	t.Parallel()

	qi := []string{"x", "y"}
	s := []string{"original", "z"}

	tree := sigo.NewKDTree(3, 2, 2, make(map[string]int))
	node := sigo.NewNode(&tree, "root", 0)

	record1 := createRowReidentification(5, 6, qi, json.Number("1"), s, "")
	node.Add(record1)

	record2 := createRowReidentification(7, 6.67, qi, json.Number("0"), s, "a")
	node.Add(record2)
	node.Add(createRowReidentification(8, 9, qi, json.Number("0"), s, "b"))
	node.Add(createRowReidentification(8, 4, qi, json.Number("1"), s, ""))
	node.Add(createRowReidentification(5, 3, qi, json.Number("0"), s, "c"))
	node.Add(createRowReidentification(8, 10, qi, json.Number("1"), s, ""))

	anonymizer := sigo.NewReidentification()
	anonymizedRecord := anonymizer.Anonymize(record1, node.Clusters()[0], []string{"x", "y"}, []string{"original", "z"})

	log.Println(anonymizedRecord.Row())
}

func createRow(x, y float64, qi []string, z string, s []string) infra.JSONLineRecord {
	row := jsonline.NewRow()
	row.Set("x", x)
	row.Set("y", y)
	row.Set("z", z)

	return infra.NewJSONLineRecord(&row, &qi, &s)
}

func createRowReidentification(x, y float64, qi []string, o json.Number, s []string, z string) infra.JSONLineRecord {
	row := jsonline.NewRow()
	row.Set("x", x)
	row.Set("y", y)
	row.Set("z", z)
	row.Set("original", o)

	return infra.NewJSONLineRecord(&row, &qi, &s)
}

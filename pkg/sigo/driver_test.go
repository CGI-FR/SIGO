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
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestSimpleClustering(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("ID", "1")

	sourceText := `{"x":0, "y":0, "foo":"bar"}
				   {"x":1, "y":1, "foo":"bar"}
				   {"x":0, "y":1, "foo":"bar"}
				   {"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	assert.Nil(t, err)

	result := []map[string]interface{}{}
	sink := infra.NewSliceDictionariesSink(&result)
	err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), 2, 1, 2, sigo.NewNoAnonymizer(), sink, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, json.Number("0"), result[0]["x"])
	assert.Equal(t, json.Number("0"), result[0]["y"])
	assert.Equal(t, "bar", result[0]["foo"])
}

func TestClusteringInfos(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("ID", "1")
}

func TestGeneralizedClustering(t *testing.T) {
	t.Parallel()

	sourceText := `{"x":0, "y":0, "foo":"bar"}
				   {"x":1, "y":1, "foo":"bar"}
				   {"x":0, "y":1, "foo":"bar"}
				   {"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	source, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	assert.Nil(t, err)

	result := []map[string]interface{}{}
	sink := infra.NewSliceDictionariesSink(&result)
	err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), 2, 1, 2, sigo.NewNoAnonymizer(), sink,
		sigo.NewSequenceDebugger("clusterID"))
	assert.Nil(t, err)

	assert.Equal(t, 1, result[0]["clusterID"])
	assert.Equal(t, 1, result[1]["clusterID"])
	assert.Equal(t, 1, result[2]["clusterID"])
	assert.Equal(t, 2, result[3]["clusterID"])
	assert.Equal(t, 2, result[4]["clusterID"])
	assert.Equal(t, 2, result[5]["clusterID"])
}

//nolint: gochecknoglobals
var tests = []struct {
	n int
	a sigo.Anonymizer
	s string
}{
	{n: 1000, a: sigo.NewNoAnonymizer(), s: "NoAnonymizer"},
	{n: 1000, a: sigo.NewGeneralAnonymizer(), s: "Generalization"},
	{n: 1000, a: sigo.NewAggregationAnonymizer("mean"), s: "MeanAggregation"},
	{n: 1000, a: sigo.NewCodingAnonymizer(), s: "TopBottomCoding"},
	{n: 1000, a: sigo.NewNoiseAnonymizer("gaussian"), s: "GaussianNoise"},
	{n: 100000, a: sigo.NewNoAnonymizer(), s: "NoAnonymizer"},
	{n: 100000, a: sigo.NewGeneralAnonymizer(), s: "Generalization"},
	{n: 100000, a: sigo.NewAggregationAnonymizer("mean"), s: "MeanAggregation"},
	{n: 100000, a: sigo.NewCodingAnonymizer(), s: "TopBottomCoding"},
	{n: 100000, a: sigo.NewNoiseAnonymizer("gaussian"), s: "GaussianNoise"},
	{n: 1000000, a: sigo.NewNoAnonymizer(), s: "NoAnonymizer"},
	{n: 1000000, a: sigo.NewGeneralAnonymizer(), s: "Generalization"},
	{n: 1000000, a: sigo.NewAggregationAnonymizer("mean"), s: "MeanAggregation"},
	{n: 1000000, a: sigo.NewCodingAnonymizer(), s: "TopBottomCoding"},
	{n: 1000000, a: sigo.NewNoiseAnonymizer("gaussian"), s: "GaussianNoise"},
}

func BenchmarkAnonymize(b *testing.B) {
	for _, test := range tests {
		b.Run(fmt.Sprintf("input_size_%d, anonymizer_%s", test.n, test.s), func(b *testing.B) {
			sourceText := []string{}

			for i := 0; i < test.n; i++ {
				// nolint: gosec
				x := rand.Intn(test.n)
				// nolint: gosec
				y := rand.Intn(test.n)
				sourceText = append(sourceText, `{"x:"`+strconv.Itoa(x)+`, "y":`+strconv.Itoa(y)+`}`)
			}

			data := strings.Join(sourceText, "\n")

			source, _ := infra.NewJSONLineSource(strings.NewReader(data), []string{"x", "y"}, []string{"foo"})

			result := []map[string]interface{}{}
			sink := infra.NewSliceDictionariesSink(&result)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = sigo.Anonymize(source, sigo.NewKDTreeFactory(), 3, 1, 2, test.a, sink, sigo.NewNoDebugger())
			}
		})
	}
}

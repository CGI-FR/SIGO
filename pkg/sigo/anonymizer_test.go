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
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/stretchr/testify/assert"
)

func TestMicroAggregationAnonymizer(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("ID", "1")

	//nolint: goconst
	sourceText := `{"x":0, "y":0, "foo":"bar"}
				   {"x":1, "y":1, "foo":"bar"}
				   {"x":0, "y":1, "foo":"bar"}
				   {"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	source1, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	assert.Nil(t, err)

	resultMean := []map[string]interface{}{}
	sink1 := infra.NewSliceDictionariesSink(&resultMean)
	err = sigo.Anonymize(source1, sigo.NewKDTreeFactory(), 2, 1, 2,
		sigo.NewAggregationAnonymizer("mean"), sink1, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, 0.33, resultMean[0]["x"])
	assert.Equal(t, 0.67, resultMean[0]["y"])
	assert.Equal(t, 2.33, resultMean[3]["x"])
	assert.Equal(t, 2.00, resultMean[3]["y"])
	assert.Equal(t, "bar", resultMean[0]["foo"])
	assert.Equal(t, "baz", resultMean[3]["foo"])

	source2, err := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	assert.Nil(t, err)

	resultMedian := []map[string]interface{}{}
	sink2 := infra.NewSliceDictionariesSink(&resultMedian)
	err = sigo.Anonymize(source2, sigo.NewKDTreeFactory(), 2, 1, 2,
		sigo.NewAggregationAnonymizer("median"), sink2, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, float64(0), resultMedian[0]["x"])
	assert.Equal(t, float64(1), resultMedian[0]["y"])
	assert.Equal(t, float64(2), resultMedian[3]["x"])
	assert.Equal(t, float64(2), resultMedian[3]["y"])
	assert.Equal(t, "bar", resultMedian[0]["foo"])
	assert.Equal(t, "baz", resultMedian[3]["foo"])
}

func TestTopBottomCodingAnonymizer(t *testing.T) {
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
	err = sigo.Anonymize(source, sigo.NewKDTreeFactory(), 7, 1, 2, sigo.NewCodingAnonymizer(),
		sink, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, 1.00, result[0]["y"])
	assert.Equal(t, 1.00, result[6]["x"])
	assert.Equal(t, 12.00, result[5]["y"])
	assert.Equal(t, 5.50, result[7]["x"])
	assert.Equal(t, 5.50, result[8]["x"])
	assert.Equal(t, 5.00, result[8]["y"])
	assert.Equal(t, 11.50, result[11]["y"])
	assert.Equal(t, 11.50, result[14]["y"])
	assert.Equal(t, 11.50, result[14]["x"])
	assert.Equal(t, 11.50, result[13]["x"])
}

func TestRandomNoiseAnonymizer(t *testing.T) {
	t.Parallel()

	tests := []int{10, 20, 50, 100, 200, 500, 1000}

	// nolint: paralleltest
	for i, N := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			t.Parallel()

			rand.Seed(time.Now().UnixNano())

			kdtree := sigo.NewKDTreeFactory().New(3, 1, 1)
			rows := []jsonline.Row{}

			for i := 0; i < N; i++ {
				row := jsonline.NewRow()
				//nolint: gosec
				x := float64(rand.Intn(N)) + rand.Float64()
				row.Set("x", x)
				rows = append(rows, row)
			}

			for i := 0; i < N; i++ {
				record := infra.NewJSONLineRecord(&rows[i], &[]string{"x"}, &[]string{})

				kdtree.Add(record)
			}

			kdtree.Build()
		})
	}
}

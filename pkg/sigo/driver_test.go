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
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestSimpleClustering(t *testing.T) {
	t.Parallel()

	//nolint: goconst
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

type LoopReader struct {
	template string
	n        int
}

func (lr *LoopReader) Read(b []byte) (int, error) {
	n, err := strings.NewReader(lr.template).Read(b)
	fmt.Println(len(b))
	// fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
	fmt.Printf("b[:n] = %q\n", b[:n])

	if err == io.EOF {
		if lr.n == 0 {
			return n, io.EOF
		}

		fmt.Println(lr.n)

		lr.n--

		return n, nil
	}

	//nolint: wrapcheck
	return n, err
}

func BenchmarkSimpleClustering(b *testing.B) {
	iter := `{"x":0, "y":0, "foo":"bar"}
				   {"x":1, "y":1, "foo":"bar"}
				   {"x":0, "y":1, "foo":"bar"}
				   {"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	for i := 0; i < b.N; i++ {
		source, err := infra.NewJSONLineSource(strings.NewReader(iter), []string{"x", "y"}, []string{"foo"})
		assert.Nil(b, err)
		b.StartTimer()

		err = sigo.Anonymize(
			source, sigo.NewKDTreeFactory(),
			2, 1, 2,
			sigo.NewAggregationAnonymizer("mean"),
			infra.NewJSONLineSink(io.Discard), sigo.NewNoDebugger(),
		)

		assert.Nil(b, err)
		b.StopTimer()
	}
}

func BenchmarkLongClustering(b *testing.B) {
	for i := 0; i < b.N; i++ {
		jsonBytes, err := os.Open("testdata/tree.json")

		assert.Nil(b, err)

		source, err := infra.NewJSONLineSource(jsonBytes, []string{"x", "y"}, []string{})
		assert.Nil(b, err)
		b.StartTimer()

		err = sigo.Anonymize(
			source,
			sigo.NewKDTreeFactory(),
			10, 1, 2,
			sigo.NewAggregationAnonymizer("mean"),
			infra.NewJSONLineSink(io.Discard), sigo.NewNoDebugger(),
		)

		assert.Nil(b, err)
		b.StopTimer()

		jsonBytes.Close()
	}
}

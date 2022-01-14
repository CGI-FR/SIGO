// Copyright (C) 2021 CGI France
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

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestMicroAggregationAnonymizer(t *testing.T) {
	t.Parallel()

	row := jsonline.NewRow()
	row.Set("ID", "1")

	sourceText := `{"x":0, "y":0, "foo":"bar"}
				   {"x":1, "y":1, "foo":"bar"}
				   {"x":0, "y":1, "foo":"bar"}
				   {"x":2, "y":1, "foo":"baz"}
				   {"x":3, "y":2, "foo":"baz"}
				   {"x":2, "y":3, "foo":"baz"}`

	source1 := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	resultMean := []map[string]interface{}{}
	sink1 := infra.NewSliceDictionariesSink(&resultMean)
	err := sigo.Anonymize(source1, sigo.NewKDTreeFactory(), 2, 1, 2, sigo.NewAggregationAnonymizer("mean"), sink1, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, 0.33, resultMean[0]["x"])
	assert.Equal(t, 0.67, resultMean[0]["y"])
	assert.Equal(t, 2.33, resultMean[3]["x"])
	assert.Equal(t, 2.00, resultMean[3]["y"])
	assert.Equal(t, "bar", resultMean[0]["foo"])
	assert.Equal(t, "baz", resultMean[3]["foo"])

	source2 := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"}, []string{"foo"})
	resultMedian := []map[string]interface{}{}
	sink2 := infra.NewSliceDictionariesSink(&resultMedian)
	err = sigo.Anonymize(source2, sigo.NewKDTreeFactory(), 2, 1, 2, sigo.NewAggregationAnonymizer("median"), sink2, sigo.NewNoDebugger())
	assert.Nil(t, err)

	assert.Equal(t, float64(0), resultMedian[0]["x"])
	assert.Equal(t, float64(1), resultMedian[0]["y"])
	assert.Equal(t, float64(2), resultMedian[3]["x"])
	assert.Equal(t, float64(2), resultMedian[3]["y"])
	assert.Equal(t, "bar", resultMedian[0]["foo"])
	assert.Equal(t, "baz", resultMedian[3]["foo"])
}

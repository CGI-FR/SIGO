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

package reidentification_test

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine", 0.5)

	row := make(map[string]interface{})
	row["x"] = 20
	row["y"] = 18

	original, err := infra.NewJSONLineSource(strings.NewReader(`{"x":20,"y":18}`), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	maskedDataset, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.InitData(original, masked)

	id.ScaleData("original", []string{"x", "y"}, []string{"z"})
	id.ScaleData("filtered", []string{"x", "y"}, []string{"z"})

	identified := id.Identify(row, row, []string{"x", "y"}, []string{"z"})

	expected := jsonline.NewRow()
	expected.Set("x", 20)
	expected.Set("y", 18)
	expected.Set("sensitive", []string{"b"})
	expected.Set("similarity", 99.18)
	recordExpected := infra.NewJSONLineRecord(&expected, &[]string{"x", "y"}, &[]string{"sensitive"})

	assert.Equal(t, recordExpected.Row(), identified.Record().Row())
}

func BenchmarkIdentify(b *testing.B) {
	id := reidentification.NewIdentifier("cosine", 0.5)

	row := make(map[string]interface{})
	row["x"] = 20
	row["y"] = 18

	data1, err := os.Open("testdata/simple.json")
	assert.Nil(b, err)

	data2, err := os.Open("testdata/maskedsimple.json")
	assert.Nil(b, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(data1), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(data2), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	id.InitData(original, masked)

	id.ScaleData("original", []string{"x", "y"}, []string{"z"})
	id.ScaleData("filtered", []string{"x", "y"}, []string{"z"})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Identify ran 100 times and each call took an average of 22357170 nanoseconds to complete.
		// 0.02 seconds for 1000 rows.
		id.Identify(row, row, []string{"x", "y"}, []string{"z"})
	}

	data1.Close()
	data2.Close()
}

func TestInitData(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine", 0.5)

	original, err := infra.NewJSONLineSource(strings.NewReader(`{"x":20,"y":18}`), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	maskedDataset, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.InitData(original, masked)

	res := id.ReturnGroup()

	var res1, res2 map[string]interface{}

	for _, record := range *res {
		if record["x"] == 3.00 && record["y"] == 7.00 {
			res1 = record
		}

		if record["x"] == 7.00 && record["y"] == 6.67 {
			res2 = record
		}
	}

	expected1 := jsonline.NewRow()
	expected1.Set("x", 3.00)
	expected1.Set("y", 7.00)
	expected1.Set("z", "")
	recordExpected1 := infra.NewJSONLineRecord(&expected1, &[]string{"x", "y"}, &[]string{"z"})

	assert.Equal(t, res1, recordExpected1.Row())

	expected2 := jsonline.NewRow()
	expected2.Set("x", 7.00)
	expected2.Set("y", 6.67)
	expected2.Set("z", "a")
	recordExpected2 := infra.NewJSONLineRecord(&expected2, &[]string{"x", "y"}, &[]string{"z"})

	assert.Equal(t, res2, recordExpected2.Row())
}

func BenchmarkInitData(b *testing.B) {
	id := reidentification.NewIdentifier("cosine", 0.5)

	data1, err := os.Open("testdata/simple.json")
	assert.Nil(b, err)

	data2, err := os.Open("testdata/maskedsimple.json")
	assert.Nil(b, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(data1), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(data2), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// InitData ran 1132 times and each call took an average of 1354702 nanoseconds to complete.
		// 0.001 seconds for 1000 rows.
		id.InitData(original, masked)
	}

	data1.Close()
	data2.Close()
}

func TestScaleData(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("euclidean", 0.5)

	original, err := infra.NewJSONLineSource(strings.NewReader(`{"x":20,"y":18}`), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	source := `{"x":10,"y":20,"z":"a"}
			   {"x":30,"y":40,"z":"b"}
			   {"x":50,"y":50,"z":"c"}`

	masked, err := infra.NewJSONLineSource(strings.NewReader(source), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.InitData(original, masked)

	id.ScaleData("filtered", []string{"x", "y"}, []string{"z"})

	res := *id.ReturnScaled("filtered")

	dataScaled := make([]map[string]interface{}, 3)
	xScaled := []float64{-1.224744871391589, 0, 1.224744871391589}
	yScaled := []float64{-1.3365734230692703, 0.26699397113501305, 1.0687776682371546}
	zScaled := []string{"a", "b", "c"}

	for i := range xScaled {
		vals := make(map[string]interface{})
		vals["x"] = xScaled[i]
		vals["y"] = yScaled[i]
		vals["z"] = zScaled[i]

		dataScaled[i] = vals
	}

	assert.Contains(t, dataScaled, res[0])
	assert.Contains(t, dataScaled, res[1])
	assert.Contains(t, dataScaled, res[2])
}

func BenchmarkScaleData(b *testing.B) {
	id := reidentification.NewIdentifier("cosine", 0.5)

	data1, err := os.Open("testdata/simple.json")
	assert.Nil(b, err)

	data2, err := os.Open("testdata/maskedsimple.json")
	assert.Nil(b, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(data1), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(data2), []string{"x", "y"}, []string{"z"})
	assert.Nil(b, err)

	id.InitData(original, masked)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Scaledata ran 1 times and each call took an average of 2783410300 nanoseconds to complete.
		// 2 seconds for 1000 rows.
		id.ScaleData("original", []string{"x", "y"}, []string{"z"})
		id.ScaleData("filtered", []string{"x", "y"}, []string{"z"})
	}

	data1.Close()
	data2.Close()
}

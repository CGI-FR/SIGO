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
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestReIdentify(t *testing.T) {
	t.Parallel()

	originalFile, err := os.Open("../../examples/re-identification/openData.json")
	assert.Nil(t, err)

	original, err := infra.NewJSONLineSource(bufio.NewReader(originalFile), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	maskedFile, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedFile), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	result := []map[string]interface{}{}
	sink := infra.NewSliceDictionariesSink(&result)

	err = reidentification.ReIdentify(original, masked, reidentification.NewIdentifier("canberra", 0.5), sink)
	assert.Nil(t, err)

	assert.Equal(t, json.Number("8"), result[0]["x"])
	assert.Equal(t, json.Number("4"), result[0]["y"])
	assert.Equal(t, []string{"a"}, result[0]["sensitive"])
}

func BenchmarkSimpleReidentification(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data1, err := os.Open("testdata/simple.json")
		assert.Nil(b, err)

		data2, err := os.Open("testdata/maskedsimple.json")
		assert.Nil(b, err)

		original, err := infra.NewJSONLineSource(bufio.NewReader(data1), []string{"x", "y"}, []string{"z"})
		assert.Nil(b, err)

		masked, err := infra.NewJSONLineSource(bufio.NewReader(data2), []string{"x", "y"}, []string{"z"})
		assert.Nil(b, err)

		b.StartTimer()

		// ReIdentify ran 1 times and each call took an average of 26056784800 nanoseconds to complete.
		// 26 seconds for 1000 rows.
		err = reidentification.ReIdentify(
			original,
			masked,
			reidentification.NewIdentifier("euclidean", 0.5),
			infra.NewJSONLineSink(io.Discard),
		)

		assert.Nil(b, err)
		b.StopTimer()

		data1.Close()

		data2.Close()
	}
}

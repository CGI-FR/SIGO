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

	err = reidentification.ReIdentify(original, masked, reidentification.NewIdentifier("canberra"), sink)
	assert.Nil(t, err)

	assert.Equal(t, json.Number("8"), result[0]["x"])
	assert.Equal(t, json.Number("4"), result[0]["y"])
	assert.Equal(t, []string{"a"}, result[0]["sensitive"])
}

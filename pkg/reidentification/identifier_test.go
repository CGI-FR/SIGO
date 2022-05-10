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
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine")

	row := jsonline.NewRow()
	row.Set("x", 20)
	row.Set("y", 18)
	original := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{"z"})

	maskedDataset, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.SaveMasked(masked)
	id.GroupMasked([]string{"x", "y"}, []string{"z"})

	identified := id.Identify(original, []string{"x", "y"}, []string{"z"})

	expected := jsonline.NewRow()
	expected.Set("x", 20)
	expected.Set("y", 18)
	expected.Set("sensitive", []string{"b"})
	recordExpected := infra.NewJSONLineRecord(&expected, &[]string{"x", "y"}, &[]string{"sensitive"})

	assert.Equal(t, recordExpected.Row(), identified.Record().Row())
}

func TestGroupAnonymizedData(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine")

	maskedDataset, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.SaveMasked(masked)
	id.GroupMasked([]string{"x", "y"}, []string{"z"})

	res := id.ReturnGroup()

	var res1, res2 map[string]interface{}

	for _, record := range *res {
		if record["x"] == "3" && record["y"] == "7" {
			res1 = record
		}

		if record["x"] == "7" && record["y"] == "6.67" {
			res2 = record
		}
	}

	expected1 := jsonline.NewRow()
	expected1.Set("x", "3")
	expected1.Set("y", "7")
	expected1.Set("z", "")
	recordExpected1 := infra.NewJSONLineRecord(&expected1, &[]string{"x", "y"}, &[]string{"z"})

	assert.Equal(t, res1, recordExpected1.Row())

	expected2 := jsonline.NewRow()
	expected2.Set("x", "7")
	expected2.Set("y", "6.67")
	expected2.Set("z", "a")
	recordExpected2 := infra.NewJSONLineRecord(&expected2, &[]string{"x", "y"}, &[]string{"z"})

	assert.Equal(t, res2, recordExpected2.Row())
}

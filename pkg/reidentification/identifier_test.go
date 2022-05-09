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
	"log"
	"os"
	"testing"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine", 3)

	row := jsonline.NewRow()
	row.Set("x", 20)
	row.Set("y", 18)
	record := infra.NewJSONLineRecord(&row, &[]string{"x", "y"}, &[]string{"z"})

	maskedDataset, err := os.Open("../../examples/re-identification/test1/data2-sigo.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.SaveMasked(masked)

	identified := id.Identify(record, masked, []string{"x", "y"}, []string{"z"})

	expected := jsonline.NewRow()
	expected.Set("x", 20)
	expected.Set("y", 18)
	expected.Set("sensitive", []string{"b"})
	recordExpected := infra.NewJSONLineRecord(&expected, &[]string{"x", "y"}, &[]string{"sensitive"})

	assert.Equal(t, recordExpected.Row(), identified.Record().Row())
}

func TestGroup(t *testing.T) {
	t.Parallel()

	id := reidentification.NewIdentifier("cosine", 3)

	maskedDataset, err := os.Open("../../examples/re-identification/anonymized.json")
	assert.Nil(t, err)

	masked, err := infra.NewJSONLineSource(bufio.NewReader(maskedDataset), []string{"x", "y"}, []string{"z"})
	assert.Nil(t, err)

	id.GroupMasked(masked, []string{"x", "y"}, []string{"z"})

	log.Println(id.ReturnsGroup())
}

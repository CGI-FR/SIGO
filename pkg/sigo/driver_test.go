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
	"encoding/json"
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
{"x":2, "y":3, "foo":"baz"}
`
	source := infra.NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"})
	result := []map[string]interface{}{}
	sink := infra.NewSliceDictionariesSink(&result)
	err := sigo.Anonymize(source, sigo.NewKDTreeFactory(), 2, 1, 2, sigo.NewNoAnonymizer(), sink)
	assert.Nil(t, err)

	assert.Equal(t, json.Number("1"), result[0]["x"])
	assert.Equal(t, json.Number("1"), result[0]["y"])
	assert.Equal(t, "bar", result[0]["foo"])
}

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
	"io"
	"strings"
	"testing"

	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func NewJSONLineRecord(row *jsonline.Row, quasiIdentifers *[]string) JSONLineRecord {
	return JSONLineRecord{row, quasiIdentifers}
}

type JSONLineRecord struct {
	row             *jsonline.Row
	quasiIdentifers *[]string
}

func (jlr JSONLineRecord) QuasiIdentifer() []float32 {
	result := []float32{}

	for _, key := range *jlr.quasiIdentifers {
		result = append(result, (*jlr.row).GetFloat32(key))
	}

	return result
}

func (jlr JSONLineRecord) Row() map[string]interface{} {
	result, err := (*jlr.row).Export()
	if err != nil {
		return nil
	}

	return result.(map[string]interface{})
}

func NewJSONLineSource(r io.Reader, quasiIdentifers []string) sigo.RecordSource {
	// nolint: exhaustivestruct
	return &JSONLineSource{importer: jsonline.NewImporter(r), quasiIdentifers: quasiIdentifers}
}

type JSONLineSource struct {
	importer        jsonline.Importer
	err             error
	record          sigo.Record
	quasiIdentifers []string
}

func (s *JSONLineSource) Err() error {
	return s.err
}

func (s *JSONLineSource) Next() bool {
	hasNext := s.importer.Import()
	if !hasNext {
		return false
	}

	row, err := s.importer.GetRow()

	s.err = err

	if s.err != nil {
		return false
	}

	s.record = NewJSONLineRecord(&row, &s.quasiIdentifers)

	return true
}

func (s *JSONLineSource) Value() sigo.Record {
	return s.record
}

func NewJSONLineSink(w io.Writer) JSONLineSink {
	return JSONLineSink{exporter: jsonline.NewExporter(w)}
}

type JSONLineSink struct {
	exporter jsonline.Exporter
}

func (s JSONLineSink) Collect(rec sigo.Record) error {
	return s.exporter.Export(rec.Row())
}

func NewSliceDictionariesSink(slice *[]map[string]interface{}) *SliceDictionariesSink {
	return &SliceDictionariesSink{slice: slice}
}

type SliceDictionariesSink struct {
	slice *[]map[string]interface{}
}

func (s *SliceDictionariesSink) Collect(rec sigo.Record) error {
	*s.slice = append(*s.slice, rec.Row())

	return nil
}

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
	source := NewJSONLineSource(strings.NewReader(sourceText), []string{"x", "y"})
	result := []map[string]interface{}{}
	sink := NewSliceDictionariesSink(&result)
	err := sigo.Anonymize(source, sigo.NewKDTreeFactory(), 2, 1, sigo.NewNoAnonymizer(), sink)
	assert.Nil(t, err)

	assert.Equal(t, json.Number("0"), result[0]["x"])
	assert.Equal(t, json.Number("0"), result[0]["y"])
	assert.Equal(t, "bar", result[0]["foo"])
}

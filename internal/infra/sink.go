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

package infra

import (
	"errors"
	"io"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

func NewJSONLineSink(w io.Writer) JSONLineSink {
	return JSONLineSink{exporter: jsonline.NewExporter(w)}
}

type JSONLineSink struct {
	exporter jsonline.Exporter
}

func (s JSONLineSink) Collect(rec sigo.Record) error {
	return errors.Unwrap(s.exporter.Export(rec.Row()))
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

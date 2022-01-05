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

package infra

import (
	"io"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

func NewJSONLineRecord(row *jsonline.Row, quasiIdentifers *[]string, sensitives *[]string) JSONLineRecord {
	return JSONLineRecord{row, quasiIdentifers, sensitives}
}

type JSONLineRecord struct {
	row             *jsonline.Row
	quasiIdentifers *[]string
	sensitives      *[]string
}

func (jlr JSONLineRecord) QuasiIdentifer() []float32 {
	result := []float32{}

	for _, key := range *jlr.quasiIdentifers {
		result = append(result, (*jlr.row).GetFloat32(key))
	}

	return result
}

func (jlr JSONLineRecord) Sensitives() []string {
	return *jlr.sensitives
}

func (jlr JSONLineRecord) Row() map[string]interface{} {
	result, err := (*jlr.row).Export()
	if err != nil {
		return nil
	}

	return result.(map[string]interface{})
}

func NewJSONLineSource(r io.Reader, quasiIdentifers []string, sensitives []string) sigo.RecordSource {
	// nolint: exhaustivestruct
	return &JSONLineSource{importer: jsonline.NewImporter(r), quasiIdentifers: quasiIdentifers, sensitives: sensitives}
}

type JSONLineSource struct {
	importer        jsonline.Importer
	err             error
	record          sigo.Record
	quasiIdentifers []string
	sensitives      []string
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

	s.record = NewJSONLineRecord(&row, &s.quasiIdentifers, &s.sensitives)

	return true
}

func (s *JSONLineSource) Value() sigo.Record {
	return s.record
}

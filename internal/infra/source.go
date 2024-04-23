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

func NewJSONLineRecord(row *jsonline.Row, quasiIdentifers *[]string, sensitives *[]string) JSONLineRecord {
	return JSONLineRecord{row, quasiIdentifers, sensitives}
}

type JSONLineRecord struct {
	row             *jsonline.Row
	quasiIdentifers *[]string
	sensitives      *[]string
}

func (jlr JSONLineRecord) QuasiIdentifer() ([]float64, error) {
	result := []float64{}

	for _, key := range *jlr.quasiIdentifers {
		value, _ := (*jlr.row).Get(key)
		if value == nil {
			//nolint: goerr113
			err := errors.New("null value in dataset")

			return []float64{}, err
		}

		result = append(result, (*jlr.row).GetFloat64(key))
	}

	return result, nil
}

func (jlr JSONLineRecord) Sensitives() []interface{} {
	result := []interface{}{}

	for _, key := range *jlr.sensitives {
		s, _ := (*jlr.row).Get(key)
		result = append(result, s)
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

func NewJSONLineSource(r io.Reader, quasiIdentifers []string, sensitives []string) (sigo.RecordSource, error) {
	// nolint: exhaustivestruct
	source := &JSONLineSource{importer: jsonline.NewImporter(r), quasiIdentifers: quasiIdentifers, sensitives: sensitives}

	//nolint: goerr113
	err := errors.New("indicate the list of quasi-identifiers")
	if len(quasiIdentifers) == 0 {
		return source, err
	}

	return source, nil
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

func (s *JSONLineSource) QuasiIdentifer() []string {
	return s.quasiIdentifers
}

func (s *JSONLineSource) Sensitive() []string {
	return s.sensitives
}

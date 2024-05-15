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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

func NewJSONLineRecord(row *jsonline.Row, quasiIdentifers *[]string, sensitives *[]string, options ...*map[string]float64) JSONLineRecord {
	if len(options) != 0 {
		return JSONLineRecord{row, quasiIdentifers, sensitives, options[0]}
	} else {
		return JSONLineRecord{row, quasiIdentifers, sensitives, &map[string]float64{}}
	}
}

type JSONLineRecord struct {
	row             *jsonline.Row
	quasiIdentifers *[]string
	sensitives      *[]string
	float64QI       *map[string]float64
}

func (jlr JSONLineRecord) QuasiIdentifer() []float64 {
	result := []float64{}

	for _, key := range *jlr.quasiIdentifers {
		result = append(result, (*jlr.row).GetFloat64(key))
	}

	return result
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

func (jlr JSONLineRecord) GetQI() map[string]float64 {
	return *jlr.float64QI
}

func NewJSONLineSource(r io.Reader, quasiIdentifers []string, sensitives []string) (sigo.RecordSource, error) {
	// nolint: exhaustivestruct
	source := &JSONLineSource{
		importer:        jsonline.NewImporter(r),
		quasiIdentifers: quasiIdentifers,
		sensitives:      sensitives,
		DataValidator:   NewFloat64DataValidator(),
	}

	//nolint: goerr113
	err := errors.New("indicate the list of quasi-identifiers")
	if len(quasiIdentifers) == 0 {
		return source, err
	}

	return source, nil
}

type JSONLineSource struct {
	sigo.DataValidator
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
		return true
	}
	rowInterface, err2 := row.Export()
	s.err = err2
	if err2 != nil {
		return true
	}

	float64Map, err3 := s.DataValidator.Validation(rowInterface.(map[string]interface{}), s.quasiIdentifers)
	s.err = err3

	if s.err != nil {
		return true
	}

	s.record = NewJSONLineRecord(&row, &s.quasiIdentifers, &s.sensitives, &float64Map)
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

type Float64DataValidator struct{}

func NewFloat64DataValidator() Float64DataValidator {
	return Float64DataValidator{}
}

func (v Float64DataValidator) Validation(row map[string]interface{}, quasiIdentifers []string) (map[string]float64, error) {
	result := make(map[string]float64)

	for _, key := range quasiIdentifers {
		// Null value check
		if row[key] == nil {
			//nolint: goerr113
			err := errors.New("null value in dataset")

			return nil, err
		}

		// Type check
		valFloat64, typeErr := transformType(row, key)
		if typeErr != nil {
			return nil, typeErr
		}

		result[key] = valFloat64
	}

	return result, nil
}

func transformType(row map[string]interface{}, key string) (float64, error) {
	var result float64

	//nolint: varnamelen
	switch t := row[key].(type) {
	case int:
		result = float64(t)
	case string:
		//nolint: gomnd
		val, err := strconv.ParseFloat(t, 64)
		if err != nil {
			//nolint: goerr113
			err = fmt.Errorf("unsupported type: %T", t)

			return result, err
		}

		result = val
	case float32:
		result = float64(t)
	case json.Number:
		val, err := t.Float64()
		if err != nil {
			//nolint: goerr113
			err = fmt.Errorf("unsupported type: %T", t)

			return result, err
		}

		result = val
	case float64:
		result = t
	default:
		//nolint: goerr113
		err := fmt.Errorf("unsupported type: %T", t)

		return result, err
	}

	return result, nil
}

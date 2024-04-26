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

package sigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type Float64DataValidator struct {
	records         []Record
	quasiIdentifers []string
}

func NewFloat64DataValidator(records []Record, quasiIdentifers []string) Float64DataValidator {
	return Float64DataValidator{records: records, quasiIdentifers: quasiIdentifers}
}

func (v Float64DataValidator) Validation() error {
	for _, record := range v.records {
		row := record.Row()

		for _, key := range v.quasiIdentifers {
			// Null value check
			if row[key] == nil {
				//nolint: goerr113
				err := errors.New("null value in dataset")

				return err
			}

			// Type check
			isValide, typeErr := checkType(row, key)
			if isValide {
				return typeErr
			}
		}
	}

	return nil
}

func checkType(row map[string]interface{}, key string) (bool, error) {
	//nolint: varnamelen
	switch t := row[key].(type) {
	case int:
		return false, nil
	case string:
		//nolint: gomnd
		_, err := strconv.ParseFloat(t, 64)
		if err != nil {
			//nolint: goerr113
			err = fmt.Errorf("unsupported type: %T", t)

			return true, err
		}
	case float32:
		return false, nil
	case json.Number:
		return false, nil
	case float64:
		return false, nil
	case []interface{}:
		return false, nil
	default:
		//nolint: goerr113
		err := fmt.Errorf("unsupported type: %T", t)

		return true, err
	}

	return false, nil
}

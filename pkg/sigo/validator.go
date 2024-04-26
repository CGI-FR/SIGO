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
	"errors"
	"fmt"
)

type Float64DataValidator struct {
	records         []Record
	quasiIdentifers []string
}

func NewFloat64DataValidator(records []Record, quasiIdentifers []string) Float64DataValidator {
	return Float64DataValidator{records: records, quasiIdentifers: quasiIdentifers}
}

// nolint: cyclop
func (v Float64DataValidator) Validation() error {
	for _, record := range v.records {
		row := record.Row()

		for _, key := range v.quasiIdentifers {
			if row[key] == nil {
				//nolint: goerr113
				err := errors.New("null value in dataset")

				return err
			}

			switch t := row[key].(type) {
			case bool:
				err := fmt.Errorf("unsupported type: %T", t)

				return err
			}
		}
	}

	return nil
}

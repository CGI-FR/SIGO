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

	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Float64DataValidator struct {
	source sigo.RecordSource
}

func NewFloat64DataValidator(source sigo.RecordSource) Float64DataValidator {
	return Float64DataValidator{source: source}
}

func (v Float64DataValidator) Validation() error {
	// Check Null value
	// if value == nil {
	// 	//nolint: goerr113
	// 	err := errors.New("null value in dataset")

	// 	return err
	// }

	// Check all data can transfer to float64
	// for _, value := range data {
	// 	_, err := strconv.ParseFloat(value, 64)
	// 	if err != nil {
	// 		return errors.New("unsupported type in dataset")
	// 	}
	// }
	//nolint: goerr113
	return errors.Unwrap(errors.New("not valide"))
}

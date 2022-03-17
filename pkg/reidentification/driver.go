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

package reidentification

import (
	"fmt"

	"github.com/cgi-fr/sigo/pkg/sigo"
)

func ReIdentify(original, masked sigo.RecordSource, identifier Identifier, sink sigo.RecordSink) error {
	for original.Next() {
		identified := identifier.Identify(original.Value(), masked, masked.QuasiIdentifer(), masked.Sensitive())

		if !identified.IsEmpty() {
			err := sink.Collect(identified.Record())
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}

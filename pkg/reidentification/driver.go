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
	"bufio"
	"fmt"
	"os"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
	"github.com/rs/zerolog/log"
)

// ReIdentify returns the list of reidentified data into sigo.RecordSink.
func ReIdentify(original, masked sigo.RecordSource, identifier Identifier, sink sigo.RecordSink) error {
	identifier.SaveData(original, "original")
	identifier.SaveData(masked, "anonymized")
	identifier.GroupMasked(masked.QuasiIdentifer(), masked.Sensitive())

	scaledOriginal := ScaleData(*identifier.original, masked.Sensitive())

	for i := range *identifier.original {
		originalValue := (*identifier.original)[i]
		originalScaledValue := scaledOriginal[i]
		identified := identifier.Identify(originalScaledValue, originalValue, masked.QuasiIdentifer(), masked.Sensitive())

		if !identified.IsEmpty() {
			err := sink.Collect(identified.Record())
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}

// LoadFile loads data into a sigo.RecordSource.
func LoadFile(file string, qi, s []string) (source sigo.RecordSource) {
	data, err := os.Open(file)
	if err != nil {
		log.Err(err).Msg("Cannot open dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	source, err = infra.NewJSONLineSource(bufio.NewReader(data), qi, s)
	if err != nil {
		log.Err(err).Msg("Cannot load jsonline dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	return source
}

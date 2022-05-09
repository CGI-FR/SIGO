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

func ReIdentify(original, masked sigo.RecordSource, identifier Identifier, sink sigo.RecordSink) error {
	identifier.SaveMasked(masked)

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

func LoadFiles(originalFile, anonymizedFile string, qi, s []string) (original, anonymized sigo.RecordSource) {
	originalData, err := os.Open(originalFile)
	if err != nil {
		log.Err(err).Msg("Cannot open original dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	original, err = infra.NewJSONLineSource(bufio.NewReader(originalData), qi, s)
	if err != nil {
		log.Err(err).Msg("Cannot load jsonline original dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	anonymizedData, err := os.Open(anonymizedFile)
	if err != nil {
		log.Err(err).Msg("Cannot open anonymized dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	anonymized, err = infra.NewJSONLineSource(bufio.NewReader(anonymizedData), qi, s)
	if err != nil {
		log.Err(err).Msg("Cannot load jsonline anonymized dataset")
		log.Warn().Int("return", 1).Msg("End SIGO")
		os.Exit(1)
	}

	return original, anonymized
}

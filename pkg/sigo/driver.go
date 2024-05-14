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
	"fmt"

	"github.com/rs/zerolog/log"
)

func Anonymize(source RecordSource, factory GeneralizerFactory,
	k int, l int, dim int, anonymyzer Anonymizer, sink RecordSink, debugger Debugger,
) error {
	generalizer := factory.New(k, l, dim, source.QuasiIdentifer())
	count := 0
	records := []Record{}

	log.Info().Msg("Reading source")

	for source.Next() {
		if source.Err() != nil {
			return fmt.Errorf("%w", source.Err())
		}

		generalizer.Add(source.Value())
		records = append(records, source.Value())
		count++
	}

	validator := NewFloat64DataValidator(records, source.QuasiIdentifer())

	valsFloat64, err := validator.Validation()
	if err != nil {
		return err
	}
	fmt.Println(valsFloat64)
	log.Info().Msgf("%v individuals to anonymize", count)
	log.Info().Msg("Tree building")

	generalizer.Build()

	log.Info().Msg("Cluster Anonymization")

	var i int64

	for _, cluster := range generalizer.Clusters() {
		log.Debug().Msgf("Cluster: %v", cluster.ID())

		for _, record := range cluster.Records() {
			anonymizedRecord := anonymyzer.Anonymize(record, cluster, source.QuasiIdentifer(), source.Sensitive())
			anonymizedRecord = debugger.Information(anonymizedRecord, cluster)

			err := sink.Collect(anonymizedRecord)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			i++
		}
	}

	log.Info().Msg("End of Anonymization")

	return nil
}

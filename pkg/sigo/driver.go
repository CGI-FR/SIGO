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

package sigo

import (
	"fmt"
)

func Anonymize(source RecordSource, factory GeneralizerFactory,
	k int, l int, dim int, anonymyzer Anonymizer, sink RecordSink, clusterID bool) error {
	generalizer := factory.New(k, l, dim)

	for source.Next() {
		if source.Err() != nil {
			return fmt.Errorf("%w", source.Err())
		}

		generalizer.Add(source.Value())
	}

	generalizer.Build(false)

	for _, cluster := range generalizer.Clusters() {
		for _, record := range cluster.Records() {
			err := sink.Collect(anonymyzer.Anonymize(record, cluster), cluster.ID(), clusterID)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}

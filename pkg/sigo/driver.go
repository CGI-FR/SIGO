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

func Anonymize(source RecordSource, factory GeneralizerFactory,
	k int, l int, anonymyzer Anonymizer, sink RecordSink) error {
	generalizer := factory.New(k, l)

	for source.Next() {
		if source.Err() != nil {
			return source.Err()
		}

		generalizer.Add(source.Value())
	}

	for _, cluster := range generalizer.Clusters() {
		for _, record := range cluster.Records() {
			err := sink.Collect(anonymyzer.Anonymize(record, cluster))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

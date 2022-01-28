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

package attack

import (
	"fmt"

	"github.com/cgi-fr/sigo/pkg/sigo"
)

func Identify(source sigo.RecordSource,
	factory sigo.GeneralizerFactory,
	dim int, sink sigo.RecordSink,
	debugger sigo.Debugger) error {
	generalizer := factory.New(1, 1, dim)

	for source.Next() {
		if source.Err() != nil {
			return fmt.Errorf("%w", source.Err())
		}

		generalizer.Add(source.Value())
	}

	generalizer.Build()

	for _, cluster := range generalizer.Clusters() {
		for _, record := range cluster.Records() {
			record = debugger.Information(record, cluster)

			err := sink.Collect(record)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}

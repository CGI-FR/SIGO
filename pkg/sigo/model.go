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

type RecordSource interface {
	Next() bool
	Err() error
	Value() Record
	QuasiIdentifer() []string
	Sensitive() []string
}

type RecordSink interface {
	Collect(Record) error
}

type Record interface {
	QuasiIdentifer() []float64
	Sensitives() []interface{}
	Row() map[string]interface{}
}

type Cluster interface {
	Records() []Record
	ID() string
	Bounds() []bounds
}

type Generalizer interface {
	Add(Record)
	Clusters() []Cluster
	String() string
	Build()
}

type GeneralizerFactory interface {
	New(k int, l int, dim int) Generalizer
}

type Anonymizer interface {
	Anonymize(Record, Cluster, []string, []string) Record
}

type Debugger interface {
	Information(Record, Cluster) Record
}

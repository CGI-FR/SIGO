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
	"sort"
)

type File struct {
	source RecordSource
	values map[string][]float64
}

func New(s RecordSource) Analyzer {
	return File{source: s, values: make(map[string][]float64)}
}

func (f File) Add(r Record) {
	for i, key := range f.source.QuasiIdentifer() {
		f.values[key] = append(f.values[key], r.QuasiIdentifer()[i])
	}
}

func (f File) Values(key string) []float64 {
	return f.values[key]
}

func (f File) CountUniqueValues() map[string]int {
	uniques := make(map[string]int)

	for _, key := range f.source.QuasiIdentifer() {
		uniques[key] = Unique(f.values[key])
	}

	return uniques
}

func Order(countUnique map[string]int) (qiOrdered []string) {
	switched := make(map[int][]string)
	slice := []int{}

	for key, count := range countUnique {
		switched[count] = append(switched[count], key)

		slice = append(slice, count)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(slice)))

	for _, val := range slice {
		qiOrdered = append(qiOrdered, switched[val]...)
		delete(switched, val)
	}

	return qiOrdered
}

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

func (f File) CountUniqueValues() map[string]int {
	uniques := make(map[string]int)

	for _, key := range f.source.QuasiIdentifer() {
		uniques[key] = Unique(f.values[key])
	}

	return uniques
}

func order(countUnique map[string]int) (qiOrdered []string) {
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

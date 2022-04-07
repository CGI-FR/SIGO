package sigo

import (
	"sort"
)

type Source struct {
	qi     map[string]int
	values map[string][]float64
}

func NewAnalyzer(qi []string) Analyzer {
	dict := make(map[string]int)
	for i, key := range qi {
		dict[key] = i
	}

	return Source{qi: dict, values: make(map[string][]float64)}
}

func (s Source) Add(r Record) {
	for key, i := range s.qi {
		s.values[key] = append(s.values[key], r.QuasiIdentifer()[i])
	}
}

func (s Source) QI(i int) string {
	return s.Order()[i]
}

func (s Source) CountUniqueValues() map[string]int {
	uniques := make(map[string]int)

	for key := range s.qi {
		uniques[key] = Unique(s.values[key])
	}

	return uniques
}

func (s Source) Order() map[int]string {
	order := make(map[int]string)
	switched := make(map[int][]string)
	slice := []int{}

	for key, count := range s.CountUniqueValues() {
		switched[count] = append(switched[count], key)

		slice = append(slice, count)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(slice)))

	i := 0

	for _, count := range slice {
		for _, qi := range switched[count] {
			order[i] = qi
			i++
		}

		delete(switched, count)
	}

	return order
}

func (s Source) Dimension(rot int) int {
	return s.qi[s.Order()[rot]]
}

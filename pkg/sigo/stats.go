package sigo

import (
	"math"
	"sort"
)

func Mean(listValues []float64) (m float64) {
	for _, val := range listValues {
		m += val
	}

	m /= float64(len(listValues))

	//nolint: gomnd
	return math.Round(m*100) / 100
}

func Median(listValues []float64) (m float64) {
	sort.Float64s(listValues)
	lenList := len(listValues)

	if lenList%2 == 0 {
		//nolint: gomnd
		return (listValues[lenList/2] + listValues[lenList/2-1]) / 2
	}

	return listValues[(lenList-1)/2]
}

type Quartiles struct {
	Q1 float64
	Q2 float64
	Q3 float64
}

func Quartile(values []float64) Quartiles {
	lenValues := len(values)
	if lenValues == 0 {
		return Quartiles{Q1: 0, Q2: 0, Q3: 0}
	}

	sort.Float64s(values)

	var c1 int

	var c2 int

	//nolint: gomnd
	if lenValues%2 == 0 {
		c1 = lenValues / 2
		c2 = lenValues / 2
	} else {
		c1 = (lenValues - 1) / 2
		c2 = c1 + 1
	}

	Q1 := Median(values[:c1])
	Q2 := Median(values)
	Q3 := Median(values[c2:])

	return Quartiles{Q1, Q2, Q3}
}

func IQR(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	qs := Quartile(values)
	iqr := qs.Q3 - qs.Q1

	return iqr
}

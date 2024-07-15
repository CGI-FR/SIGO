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
	rd "crypto/rand"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"sort"
	"time"
)

// Min returns the min value of listValues.
func Min(listValues []float64) float64 {
	sort.Float64s(listValues)

	return listValues[0]
}

// Max returns the max value of listValues.
func Max(listValues []float64) float64 {
	sort.Float64s(listValues)

	return listValues[len(listValues)-1]
}

// Mean returns the mean value of listValues.
func Mean(listValues []float64) (m float64) {
	for _, val := range listValues {
		m += val
	}

	m /= float64(len(listValues))

	//nolint: gomnd
	return math.Round(m*100) / 100
}

// Median returns the median value of listValues.
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

// Quartile returns the 1st quartile, 2nd quartile (median) and 3rd quartile of values.
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

// IQR returns the interquartile range of values (Q3-Q1).
func IQR(values []float64) float64 {
	if len(values) == 0 {
		return math.NaN()
	}

	qs := Quartile(values)
	iqr := qs.Q3 - qs.Q1

	return iqr
}

// Std returns the standard deviation of listValues.
// (https://fr.wikipedia.org/wiki/%C3%89cart_type).
func Std(listValues []float64) (s float64) {
	for _, val := range listValues {
		//nolint: gomnd
		s += math.Pow((val - Mean(listValues)), 2)
	}

	return math.Sqrt(s / float64(len(listValues)))
}

// Sum returns the sum of the elements of listValues.
func Sum(listValues []float64) (sum float64) {
	for _, val := range listValues {
		sum += val
	}

	return sum
}

// ExpNumber generates safe exponential random value.
// (https://dzone.com/articles/generating-laplace-distributed-random-values).
func ExpNumber(mean float64) float64 {
	// Returns a secure random float in [0,1).
	random, err := RandFloat()
	if err != nil {
		return 0
	}

	return -mean * math.Log(random)
}

// LaplaceNumber generates Laplace value.
// (https://dzone.com/articles/generating-laplace-distributed-random-values).
func LaplaceNumber() float64 {
	e1 := ExpNumber(1)
	e2 := ExpNumber(1)

	return e1 - e2
}

// GaussianNumber generates safe Gaussian value.
func GaussianNumber(loc float64, scale float64) float64 {
	rand.Seed(time.Now().Unix())

	// Returns 2 normal values
	z1, z2 := BoxMuller()
	numbers := []float64{z1, z2}
	//nolint: gomnd
	idx, _ := rd.Int(rd.Reader, big.NewInt(2))
	// Choose randomly one of the 2 numbers with a safe draw.
	random := numbers[idx.Int64()]

	// if normal values (scale = 1, loc = 0)
	return random*scale + loc
}

// Scaling returns the scaled value to range [-2;2] for Laplace number and [-1;1] for Gaussian Number.
// (https://en.wikipedia.org/wiki/Normalization_(statistics)).
// If no method is mentioned, standardization is used.
func Scaling(value float64, listValues []float64, method string) float64 {
	scope := Max(listValues) - Min(listValues)
	//nolint: gomnd
	switch method {
	case laplace:
		if scope == 0 {
			return -2
		}
		// Normalization Min-Max
		return -2 + ((value-Min(listValues))*4)/(scope)
	case gaussian:
		if scope == 0 {
			return -1
		}
		// Normalization Min-Max
		return -1 + ((value-Min(listValues))*2)/(scope)
	}

	// Standardization
	return (value - Mean(listValues)) / Std(listValues)
}

// Rescaling returns the resize value.
func Rescaling(value float64, listValues []float64, method string) (rescale float64) {
	//nolint: gomnd
	switch method {
	case laplace:
		rescale = Min(listValues) + ((value+2)*(Max(listValues)-Min(listValues)))/4
	case gaussian:
		rescale = Min(listValues) + ((value+1)*(Max(listValues)-Min(listValues)))/2
	default:
		rescale = value*Std(listValues) + Mean(listValues)
	}

	return rescale
}

// Returns a secure random integer in [0,max).
func RandInt(max int64) (int, error) {
	// Int (crypto/rand) returns a uniform random value in [0, max)
	val, err := rd.Int(rd.Reader, big.NewInt(max))
	//nolint: goerr113
	if err != nil {
		return 0, errors.New("cannot generate random value")
	}

	return int(val.Uint64()), nil
}

// Returns a secure random float in [0,1).
func RandFloat() (float64, error) {
	//nolint: gomnd
	val, err := RandInt(int64(math.Pow10(15)))
	if err != nil {
		return 0, err
	}

	random := float64(val) * math.Pow10(-15)

	return random, nil
}

// BoxMuller generates pairs of independent random numbers following a normal distribution.
// (https://en.wikipedia.org/wiki/Box%E2%80%93Muller_transform).
func BoxMuller() (float64, float64) {
	x, _ := RandFloat()
	y, _ := RandFloat()

	z1 := math.Sqrt(-2.0*math.Log(x)) * math.Cos(2.0*math.Pi*y)
	z2 := math.Sqrt(-2.0*math.Log(x)) * math.Sin(2.0*math.Pi*y)

	return z1, z2
}

func Unique(values []float64) int {
	tmp := make(map[float64]int)

	for _, val := range values {
		tmp[val]++
	}

	return len(tmp)
}

// Secure shuffle of the order of the elements.
func Shuffle(s []float64) []float64 {
	slice := s
	for i := range slice {
		j, err := RandInt(int64(i + 1))
		if err != nil {
			return nil
		}

		// swaps the elements with indexes i and j
		slice[i], slice[j] = slice[j], slice[i]
	}

	return slice
}

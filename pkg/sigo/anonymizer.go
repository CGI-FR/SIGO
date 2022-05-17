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
	"encoding/json"
)

const (
	laplace  = "laplace"
	gaussian = "gaussian"
)

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

func NewGeneralAnonymizer() GeneralAnonymizer {
	return GeneralAnonymizer{groupMap: make(map[Cluster]map[string]string)}
}

func NewAggregationAnonymizer(typeAgg string) AggregationAnonymizer {
	return AggregationAnonymizer{typeAggregation: typeAgg, values: make(map[string]map[string]float64)}
}

func NewCodingAnonymizer() CodingAnonymizer {
	return CodingAnonymizer{}
}

func NewNoiseAnonymizer(mechanism string) NoiseAnonymizer {
	return NoiseAnonymizer{typeNoise: mechanism}
}

func NewSwapAnonymizer() SwapAnonymizer {
	return SwapAnonymizer{swapValues: make(map[string]map[string][]float64)}
}

func NewReidentification() Reidentification {
	return Reidentification{
		masked:    make(map[string][]map[string]interface{}),
		unique:    make(map[string]bool),
		sensitive: make(map[string]interface{}),
	}
}

type (
	NoAnonymizer      struct{}
	GeneralAnonymizer struct {
		groupMap map[Cluster]map[string]string
	}
	AggregationAnonymizer struct {
		typeAggregation string
		values          map[string]map[string]float64
	}
	CodingAnonymizer struct{}
	NoiseAnonymizer  struct {
		typeNoise string
	}
	SwapAnonymizer struct {
		swapValues map[string]map[string][]float64
	}
	AnonymizedRecord struct {
		original Record
		mask     map[string]interface{}
	}

	Reidentification struct {
		masked    map[string][]map[string]interface{}
		unique    map[string]bool
		sensitive map[string]interface{}
	}
)

func (ar AnonymizedRecord) QuasiIdentifer() []float64 {
	return ar.original.QuasiIdentifer()
}

func (ar AnonymizedRecord) Sensitives() []interface{} {
	return ar.original.Sensitives()
}

func (ar AnonymizedRecord) Row() map[string]interface{} {
	original := ar.original.Row()
	for k, v := range ar.mask {
		original[k] = v
	}

	return original
}

// Anonymize returns the original record, there is no anonymization.
func (a NoAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}
	for _, q := range qi {
		mask[q] = rec.Row()[q]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// Anonymize returns the record anonymize with the method general
// the record takes the bounds of the cluster.
func (a GeneralAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	b := clus.Bounds()

	mask := map[string]interface{}{}
	for i, q := range qi {
		mask[q] = []float64{b[i].down, b[i].up}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// Anonymize returns the record anonymized with the method meanAggregarion or medianAggregation
// the record takes the aggregated values of the cluster.
func (a AggregationAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	if a.values[clus.ID()] == nil {
		a.ComputeAggregation(clus, qi)
	}

	for _, key := range qi {
		mask[key] = a.values[clus.ID()][key]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// ComputeAggregation calculates the mean (method meanAggreagtion)
// or median (method medianAggregation) value of the cluster for each qi.
func (a AggregationAnonymizer) ComputeAggregation(clus Cluster, qi []string) {
	values := listValues(clus, qi)

	valAggregation := make(map[string]float64)

	for _, key := range qi {
		switch a.typeAggregation {
		case "mean":
			valAggregation[key] = Mean(values[key])
		case "median":
			valAggregation[key] = Median(values[key])
		}
	}

	a.values[clus.ID()] = valAggregation
}

// Anonymize returns the record anonymized with the method outlier
// if the record is in the interval [Q1;Q3] then we don't change its value
// if the record is > Q3 then it takes the Q3 value
// if the record is < Q1 then it takes the Q1 value.
func (a CodingAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values := listValues(clus, qi)
	mask := map[string]interface{}{}

	for i, key := range qi {
		vals := values[key]
		q := Quartile(vals)
		bottom := q.Q1
		top := q.Q3

		val := rec.QuasiIdentifer()[i]

		switch {
		case val < bottom:
			mask[key] = bottom
		case val > top:
			mask[key] = top
		default:
			mask[key] = val
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// Anonymize returns the record anonymized with the method laplaceNoise or gaussianNoise
// the record takes as value the original value added to a Laplacian or Gaussian noise
// the anonymized value stays within the bounds of the cluster.
func (a NoiseAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values := listValues(clus, qi)
	mask := map[string]interface{}{}

	for i, key := range qi {
		val := rec.QuasiIdentifer()[i]

		laplaceVal := Scaling(val, values[key], laplace)
		gaussianVal := Scaling(val, values[key], gaussian)

		var randomVal float64

		for {
			switch a.typeNoise {
			case laplace:
				randomVal = Rescaling(laplaceVal+LaplaceNumber(), values[key], laplace)
			case gaussian:
				randomVal = Rescaling(gaussianVal+GaussianNumber(0, 1), values[key], gaussian)
			}

			if (randomVal > Min(values[key]) && randomVal < Max(values[key])) || Min(values[key]) == Max(values[key]) {
				break
			}
		}

		mask[key] = randomVal
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a SwapAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	// cluster value swapping
	if a.swapValues[clus.ID()] == nil {
		a.Swap(clus, qi)
	}

	var idx int

	// retrieve the position (idx) of the record in the cluster
	for i, r := range clus.Records() {
		if rec == r {
			idx = i
		}
	}

	for _, key := range qi {
		// retrieve the swapped value
		mask[key] = a.swapValues[clus.ID()][key][idx]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a SwapAnonymizer) Swap(clus Cluster, qi []string) {
	// retrieve the cluster values for each qi
	values := listValues(clus, qi)
	swapVal := make(map[string][]float64)

	for _, key := range qi {
		// values permutation
		swapVal[key] = Shuffle(values[key])
	}

	a.swapValues[clus.ID()] = swapVal
}

func (r Reidentification) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}

	// initialize re-identification object
	if r.masked[clus.ID()] == nil {
		r.InitReidentification(clus, qi)
	}

	original, _ := rec.Sensitives()[0].(json.Number).Int64()

	if original == 1 {
		for _, q := range qi {
			mask[q] = rec.Row()[q]
		}

		if r.sensitive[clus.ID()] != nil {
			mask[s[1]] = r.sensitive[clus.ID()]
		} else if !r.unique[clus.ID()] {
			scores := r.ComputeSimilarity(rec, clus, qi, s)
			_, sens := TopSimilarity(scores)
			mask[s[1]] = sens
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// InitReidentification initialize the re-identification object.
func (r Reidentification) InitReidentification(clus Cluster, qi []string) {
	maskedData := []map[string]interface{}{}

	sens := []string{}

	for _, rec := range clus.Records() {
		original, _ := rec.Sensitives()[0].(json.Number).Int64()
		if original == 0 {
			maskedData = append(maskedData, rec.Row())
			sens = append(sens, rec.Sensitives()[1].(string))
		}
	}

	// groups all the anonymized records
	r.masked[clus.ID()] = maskedData

	// indicates if the cluster contains unique masked data
	r.unique[clus.ID()] = Unique(maskedData, qi)

	// checks if the sensitive data is well represented
	if IsUnique(sens) {
		r.sensitive[clus.ID()] = sens[0]
	} else {
		r.sensitive[clus.ID()] = nil
	}
}

func (r Reidentification) ComputeSimilarity(rec Record, clus Cluster,
	qi []string, s []string) map[float64]interface{} {
	scores := make(map[float64]interface{})

	x := make(map[string]interface{})
	for _, q := range qi {
		x[q] = rec.Row()[q]
	}

	X := MapItoMapF(x)

	for _, row := range r.masked[clus.ID()] {
		y := make(map[string]interface{})
		for _, q := range qi {
			y[q] = row[q]
		}

		Y := MapItoMapF(y)

		score := Similarity(ComputeDistance("", X, Y))

		scores[score] = row[s[1]]
	}

	return scores
}

// Returns the list of values present in the cluster for each qi.
func listValues(clus Cluster, qi []string) (mapValues map[string][]float64) {
	mapValues = make(map[string][]float64)

	for _, record := range clus.Records() {
		for i, key := range qi {
			mapValues[key] = append(mapValues[key], record.QuasiIdentifer()[i])
		}
	}

	return mapValues
}

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

const (
	laplace  = "laplace"
	gaussian = "gaussian"
)

type bounds struct {
	down, up float64
}

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

func NewGeneralAnonymizer() GeneralAnonymizer {
	return GeneralAnonymizer{boundsValues: make(map[string]map[string]bounds)}
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

type (
	NoAnonymizer      struct{}
	GeneralAnonymizer struct {
		// groupMap map[Cluster]map[string]string
		// map of cluster -> qi -> bounds
		boundsValues map[string]map[string]bounds
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
	mask := map[string]interface{}{}

	if a.boundsValues[clus.ID()] == nil {
		a.ComputeGeneralization(clus, qi)
	}

	for _, key := range qi {
		mask[key] = []float64{a.boundsValues[clus.ID()][key].down, a.boundsValues[clus.ID()][key].up}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

// ComputeGeneralization calculates the min and max values of the cluster for each qi.
func (a GeneralAnonymizer) ComputeGeneralization(clus Cluster, qi []string) {
	values := listValues(clus, qi)

	boundsVal := make(map[string]bounds)

	for _, key := range qi {
		var b bounds
		b.down = Min(values[key])
		b.up = Max(values[key])
		boundsVal[key] = b
	}

	a.boundsValues[clus.ID()] = boundsVal
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

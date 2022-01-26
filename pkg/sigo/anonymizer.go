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

func (a NoAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	mask := map[string]interface{}{}
	for _, q := range qi {
		mask[q] = rec.Row()[q]
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a GeneralAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	b := clus.Bounds()

	mask := map[string]interface{}{}
	for i, q := range qi {
		mask[q] = []float64{b[i].down, b[i].up}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

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

func listValues(clus Cluster, qi []string) (mapValues map[string][]float64) {
	mapValues = make(map[string][]float64)

	for _, record := range clus.Records() {
		for i, key := range qi {
			mapValues[key] = append(mapValues[key], record.QuasiIdentifer()[i])
		}
	}

	return mapValues
}

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
	"errors"
)

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

func NewGeneralAnonymizer() GeneralAnonymizer {
	return GeneralAnonymizer{groupMap: make(map[Cluster]map[string]string)}
}

func NewAggregationAnonymizer(typeAgg string) AggregationAnonymizer {
	return AggregationAnonymizer{typeAggregation: typeAgg}
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

func (ar AnonymizedRecord) QuasiIdentifer() []float32 {
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
		mask[q] = []float32{b[i].down, b[i].up}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a AggregationAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values := make(map[string][]float64)

	for _, record := range clus.Records() {
		for key, value := range record.Row() {
			switch v := value.(type) {
			case json.Number:
				val, _ := v.Float64()
				values[key] = append(values[key], val)
			case int:
				values[key] = append(values[key], float64(v))
			default:
				continue
			}
		}
	}

	mask := map[string]interface{}{}

	for _, key := range qi {
		switch a.typeAggregation {
		case "mean":
			mask[key] = Mean(values[key])
		case "median":
			mask[key] = Median(values[key])
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a CodingAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
	values := listValues(clus)
	mask := map[string]interface{}{}

	for _, key := range qi {
		vals := values[key]
		q := Quartile(vals)
		bottom := q.Q1
		top := q.Q3

		val, err := convertToFloat64(rec.Row()[key])
		if err != nil {
			continue
		}

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
	mask := map[string]interface{}{}

	for _, key := range qi {
		val, err := convertToFloat64(rec.Row()[key])
		if err != nil {
			continue
		}

		switch a.typeNoise {
		case "laplace":
			mask[key] = val + LaplaceNumber()
		case "gaussian":
			mask[key] = val + GaussianNumber(0, 1)
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func listValues(clus Cluster) (mapValues map[string][]float64) {
	mapValues = make(map[string][]float64)

	for _, record := range clus.Records() {
		for key, value := range record.Row() {
			switch v := value.(type) {
			case json.Number:
				val, _ := v.Float64()
				mapValues[key] = append(mapValues[key], val)
			case int:
				mapValues[key] = append(mapValues[key], float64(v))
			default:
				continue
			}
		}
	}

	return mapValues
}

func convertToFloat64(value interface{}) (val float64, err error) {
	if value == nil {
		//nolint: goerr113
		return 0, errors.New("error: in conversion to Float64")
	}

	switch v := value.(type) {
	case string:
		//nolint: goerr113
		return 0, errors.New("error: in conversion to Float64")
	case json.Number:
		val, _ = v.Float64()
	case int:
		val = float64(v)
	case interface{}:
		return v.(float64), nil
	}

	return val, nil
}

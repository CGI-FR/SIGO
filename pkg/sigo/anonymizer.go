package sigo

import (
	"encoding/json"
	"math"
	"math/rand"
	"sort"
)

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

func NewAggregationAnonymizer(typeAgg string) AggregationAnonymizer {
	return AggregationAnonymizer{typeAggregation: typeAgg}
}

type (
	NoAnonymizer          struct{}
	AggregationAnonymizer struct {
		typeAggregation string
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

func (a NoAnonymizer) Anonymize(rec Record, clus Cluster) Record {
	//nolint: gosec
	choice := clus.Records()[rand.Intn(len(clus.Records()))]

	for {
		if choice != rec || len(clus.Records()) < 2 {
			break
		}
		//nolint: gosec
		choice = clus.Records()[rand.Intn(len(clus.Records()))]
	}

	mask := map[string]interface{}{}
	mask["x"] = choice.Row()["x"]
	mask["y"] = choice.Row()["y"]

	return AnonymizedRecord{original: rec, mask: mask}
}

func (a AggregationAnonymizer) Anonymize(rec Record, clus Cluster) Record {
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

	for key := range values {
		switch a.typeAggregation {
		case "mean":
			mask[key] = mean(values[key])
		case "median":
			mask[key] = median(values[key])
		}
	}

	return AnonymizedRecord{original: rec, mask: mask}
}

func mean(listValues []float64) (m float64) {
	for _, val := range listValues {
		m += val
	}

	m /= float64(len(listValues))
	//nolint: gomnd
	return math.Round(m*100) / 100
}

func median(listValues []float64) (m float64) {
	sort.Float64s(listValues)
	lenList := len(listValues)

	if lenList%2 == 0 {
		//nolint: gomnd
		return (listValues[lenList/2] + listValues[lenList/2-1]) / 2
	}

	return listValues[(lenList-1)/2]
}

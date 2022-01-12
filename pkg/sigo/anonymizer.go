package sigo

import (
	"math/rand"
)

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

type NoAnonymizer struct{}

func NewGeneralAnonymizer() GeneralAnonymizer {
	return GeneralAnonymizer{groupMap: make(map[Cluster]map[string]string)}
}

type (
	GeneralAnonymizer struct {
		groupMap map[Cluster]map[string]string
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

func (a GeneralAnonymizer) Anonymize(rec Record, clus Cluster) Record {
	b := clus.Bounds()

	mask := map[string]interface{}{}
	mask["x"] = []float32{b[0].down, b[0].up}
	mask["y"] = []float32{b[1].down, b[1].up}

	return AnonymizedRecord{original: rec, mask: mask}
}

func MinMax(array []float32) (float32, float32) {
	max := array[0]
	min := array[0]

	for _, value := range array {
		if max < value {
			max = value
		}

		if min > value {
			min = value
		}
	}

	return min, max
}

package sigo

import (
	"math/rand"
)

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

type NoAnonymizer struct{}

type AnonymizedRecord struct {
	original Record
	mask     map[string]interface{}
}

func (ar AnonymizedRecord) QuasiIdentifer() []float32 {
	return ar.original.QuasiIdentifer()
}

func (ar AnonymizedRecord) Sensitives() []string {
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

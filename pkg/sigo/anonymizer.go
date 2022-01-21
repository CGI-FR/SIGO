package sigo

import (
	"math/rand"
)

type Anonymizers struct {
	dict map[string]Anonymizer
}

func NewAnonymizers() Anonymizers {
	dict := make(map[string]Anonymizer)
	dict["NoAnonymizer"] = NewNoAnonymizer()
	dict["general"] = NewGeneralAnonymizer()

	return Anonymizers{dict: dict}
}

func (anonymizers *Anonymizers) Anonymizer(name string) Anonymizer {
	return anonymizers.dict[name]
}

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

func (a NoAnonymizer) Anonymize(rec Record, clus Cluster, qi, s []string) Record {
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
	for _, q := range qi {
		mask[q] = choice.Row()[q]
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

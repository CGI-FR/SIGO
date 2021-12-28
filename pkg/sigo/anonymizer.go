package sigo

func NewNoAnonymizer() NoAnonymizer { return NoAnonymizer{} }

type NoAnonymizer struct{}

func (a NoAnonymizer) Anonymize(rec Record, clus Cluster) Record {
	return rec
}

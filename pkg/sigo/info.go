package sigo

func NewSequenceDebugger() SequenceDebugger { return SequenceDebugger{map[string]int{}} }

type SequenceDebugger struct {
	cache map[string]int
}

type InfosRecord struct {
	original Record
	infos    map[string]interface{}
}

func (ir InfosRecord) QuasiIdentifer() []float32 {
	return ir.original.QuasiIdentifer()
}

func (ir InfosRecord) Sensitives() []interface{} {
	return ir.original.Sensitives()
}

func (ir InfosRecord) Row() map[string]interface{} {
	original := ir.original.Row()
	for k, v := range ir.infos {
		original[k] = v
	}

	return original
}

func (d SequenceDebugger) id(c Cluster) int {
	count := len(d.cache)
	if d.cache[c.ID()] == 0 {
		d.cache[c.ID()] = count + 1
	}

	return d.cache[c.ID()]
}

func (d SequenceDebugger) Information(rec Record, cluster Cluster, key string) Record {
	infos := make(map[string]interface{})

	infos[key] = d.id(cluster)

	return InfosRecord{original: rec, infos: infos}
}

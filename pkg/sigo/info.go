package sigo

func NewInfos() Infos { return Infos{} }

type Infos struct{}

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

func (i Infos) Information(rec Record, cluster Cluster, key string) Record {
	infos := make(map[string]interface{})
	for k, v := range rec.Row() {
		infos[k] = v
	}

	switch key {
	case "clusterID":
		infos["clusterID"] = cluster.ClusterInfos()["ID"]
	case "clusterPath":
		infos["clusterPath"] = cluster.ClusterInfos()["Path"]
	}

	return InfosRecord{original: rec, infos: infos}
}

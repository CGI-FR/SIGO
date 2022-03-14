package reidentification

type Similarity struct {
	id        int
	individu  []float64
	score     float64
	sensitive string
}

func NewSimilarity(id int, ind []float64, sc float64, s string) Similarity {
	return Similarity{id: id, individu: ind, score: sc, sensitive: s}
}

type Individu struct {
	data []float64
	sim  []Similarity
}

func NewIndividu(data []float64, sim []Similarity) Individu {
	return Individu{data: data, sim: sim}
}

func (ind Individu) Values() []float64 {
	return ind.data
}

func (ind Individu) Similarities() []Similarity {
	return ind.sim
}

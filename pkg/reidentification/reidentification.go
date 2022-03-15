package reidentification

type Similarity struct {
	id        int
	qi        map[string]interface{}
	score     float64
	sensitive string
}

func NewSimilarity(id int, data map[string]interface{}, sc float64, s string) Similarity {
	return Similarity{id: id, qi: data, score: sc, sensitive: s}
}

type Individu struct {
	id  int
	qi  map[string]interface{}
	sim []Similarity
}

func NewIndividu(id int, data map[string]interface{}, sim []Similarity) Individu {
	return Individu{id: id, qi: data, sim: sim}
}

func (ind Individu) Values() map[string]interface{} {
	return ind.qi
}

func (ind Individu) Similarities() []Similarity {
	return ind.sim
}

func (ind Individu) Reidenfication(k int) map[string]interface{} {
	res := make(map[string]interface{})

	slice := TopSimilarity(ind.sim, k)

	if Risk(slice) == 1 {
		sensitive := Recover(slice)

		for key, val := range ind.qi {
			res[key] = val
		}

		res["sensitive"] = sensitive
	}

	return res
}

func Risk(slice []Similarity) float64 {
	sensitives := []string{}

	for _, sim := range slice {
		sensitives = append(sensitives, sim.sensitive)
	}

	count := CountValues(sensitives)
	if len(count) == 1 {
		return 1
	}

	return float64(len(count)) / float64(len(sensitives))
}

func CountValues(sensitives []string) map[string]int {
	count := make(map[string]int)
	for _, val := range sensitives {
		count[val]++
	}

	return count
}

func Recover(slice []Similarity) string {
	return slice[0].sensitive
}

package reidentification

type Original struct {
	data []Individu
}

func NewOriginal() Original {
	return Original{data: []Individu{}}
}

type Individu struct {
	id  int
	qi  map[string]interface{}
	sim []Similarity
}

func NewIndividu(id int, data map[string]interface{}, sim []Similarity) Individu {
	return Individu{id: id, qi: data, sim: sim}
}

func (or Original) Values(i int) map[string]interface{} {
	return or.data[i].qi
}

func (or Original) Similarities(i int) []Similarity {
	return or.data[i].sim
}

func (or *Original) Add(ind Individu) {
	or.data = append(or.data, ind)
}

func (or Original) Reidenfication(k int) (res []map[string]interface{}) {
	for i := range or.data {
		m := make(map[string]interface{})
		slice := TopSimilarity(or.data[i].sim, k)
		sensitive, risk := Recover(slice)

		if risk {
			for key, val := range or.data[i].qi {
				m[key] = val
			}

			m["sensitive"] = sensitive
			res = append(res, m)
		}
	}

	return res
}

func Risk(slice []Similarity) float64 {
	sensitives := []string{}

	for _, sim := range slice {
		sensitives = append(sensitives, sim.sensitive...)
	}

	count := CountValues(sensitives)

	if len(count) == 1 {
		return 1
	}

	if len(count) == len(sensitives) {
		return 1 / float64(len(sensitives))
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

func Recover(slice []Similarity) ([]string, bool) {
	if Risk(slice) == 1 {
		return slice[0].sensitive, true
	}

	return []string{""}, false
}

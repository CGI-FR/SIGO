package reidentification

import (
	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Original struct {
	data []Individu
}

func NewOriginal() Original {
	return Original{data: []Individu{}}
}

type Individu struct {
	id  int
	qi  map[string]interface{}
	sim Similarities
}

func NewIndividu(id int, data map[string]interface{}, sim Similarities) Individu {
	return Individu{id: id, qi: data, sim: sim}
}

func (or Original) Values(i int) map[string]interface{} {
	return or.data[i].qi
}

func (or Original) Similarities(i int) []Similarity {
	return or.data[i].sim.slice
}

func (or *Original) Add(ind Individu) {
	or.data = append(or.data, ind)
}

func (or Original) Identification(k int) (res []map[string]interface{}) {
	for i := range or.data {
		m := make(map[string]interface{})
		top := or.data[i].sim.TopSimilarity(k)
		sensitive, risk := Recover(top.Slice())

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

func Reidenfication(original sigo.RecordSource, masked *sigo.RecordSource,
	metric string, k int) []map[string]interface{} {
	// Liste des invidus du dataset original avec liste de similarités
	var res Original

	i := 0

	// Pour chaque indivu du dataset original
	for original.Next() {
		qi := make(map[string]interface{})

		// On récupère les valeurs QI
		for _, q := range original.QuasiIdentifer() {
			qi[q] = original.Value().Row()[q]
		}

		var dist Metric

		switch metric {
		case "cosine":
			dist = NewCosineSimilarity()
		case "manhattan":
			dist = NewManhattanDistance()
		case "canberra":
			dist = NewCanberraDistance()
		case "chebyshev":
			dist = NewChebyshevDistance()
		case "minkowski":
			//nolint: gomnd
			dist = NewMinkowskiDistance(3)
		default:
			dist = NewEuclideanDistance()
		}

		// On créé un objet Similarities qui contient la liste des similarités avec les données anonymisées
		sims := NewSimilarities(dist)
		j := 0

		// Pour chaque individu du dataset anonymisé
		for (*masked).Next() {
			// Calcul la similarité avec l'individu original
			sim := NewSimilarity(j, (*masked).Value(), (*masked).QuasiIdentifer(), (*masked).Sensitive())
			sim.ComputeSimilarity(original.Value(), original.QuasiIdentifer(), sims.Metric())

			// Ajout à la liste des Similarités
			sims.Add(sim)
			j++
		}

		ind := NewIndividu(i, qi, sims)
		i++

		// Ajout de l'individu original
		res.Add(ind)
	}

	// Calcul ré-identification sur les individu du dataset original avec en paramètre k
	return res.Identification(k)
}

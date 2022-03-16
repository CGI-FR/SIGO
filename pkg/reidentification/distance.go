package reidentification

import (
	"math"

	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Similarity struct {
	id        int
	qi        map[string]interface{}
	score     float64
	sensitive []string
}

func NewSimilarity(id int, ind sigo.Record, qid []string, s []string) Similarity {
	x := make(map[string]interface{})
	list := []string{}

	for _, q := range qid {
		x[q] = ind.Row()[q]
	}

	for i := range s {
		list = append(list, ind.Row()[s[i]].(string))
	}

	return Similarity{id: id, qi: x, score: 0, sensitive: list}
}

type Cosine struct{}

type Euclidean struct{}

func NewCosineSimilarity() Cosine { return Cosine{} }

func NewEuclideanDistance() Euclidean { return Euclidean{} }

func (cos Cosine) Compute(x, y map[string]float64) float64 {
	var dotProduct, X, Y float64

	//nolint: gomnd
	for key := range x {
		dotProduct += x[key] * y[key]
		X += math.Pow(x[key], 2)
		Y += math.Pow(y[key], 2)
	}

	return dotProduct / (math.Sqrt(X) * math.Sqrt(Y))
}

func (eu Euclidean) Compute(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		//nolint: gomnd
		sum += math.Pow((x[key] - y[key]), 2)
	}

	return math.Sqrt(sum)
}

func (sim *Similarity) ComputeSimilarity(ind sigo.Record, qid []string, metric Metric) {
	x := make(map[string]interface{})

	for _, q := range qid {
		x[q] = ind.Row()[q]
	}

	X := MapItoMapF(x)
	Y := MapItoMapF(sim.qi)

	sim.score = metric.Compute(X, Y)
}

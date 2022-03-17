package reidentification

import (
	"math"
	"sort"
)

type Cosine struct{}

type Euclidean struct{}

type Manhattan struct{}

type Canberra struct{}

type Chebyshev struct{}

type Minkowski struct {
	p float64
}

func NewCosineSimilarity() Cosine { return Cosine{} }

func NewEuclideanDistance() Euclidean { return Euclidean{} }

func NewManhattanDistance() Manhattan { return Manhattan{} }

func NewCanberraDistance() Canberra { return Canberra{} }

func NewChebyshevDistance() Chebyshev { return Chebyshev{} }

func NewMinkowskiDistance(p float64) Minkowski { return Minkowski{p: p} }

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

func (man Manhattan) Compute(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		sum += math.Abs(x[key] - y[key])
	}

	return sum
}

func (ca Canberra) Compute(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		sum += math.Abs(x[key]-y[key]) / math.Abs(x[key]+y[key])
	}

	return sum
}

func (che Chebyshev) Compute(x, y map[string]float64) float64 {
	res := make([]float64, 0, len(x))

	for key := range x {
		res = append(res, math.Abs(x[key]-y[key]))
	}

	sort.Float64s(res)

	return res[len(res)-1]
}

func (min Minkowski) Compute(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		abs := math.Abs(x[key] - y[key])
		sum += math.Pow(abs, min.p)
	}

	return math.Pow(sum, 1/min.p)
}

package reidentification

import (
	"math"
	"sort"
)

// https://towardsdatascience.com/9-distance-measures-in-data-science-918109d069fa
// https://towardsdatascience.com/how-to-decide-the-perfect-distance-metric-for-your-machine-learning-model-2fa6e5810f11

// Cosine compute the cosine distance.
func Cosine(x, y map[string]float64) float64 {
	var dotProduct, X, Y float64

	//nolint: gomnd
	for key := range x {
		dotProduct += x[key] * y[key]
		X += math.Pow(x[key], 2)
		Y += math.Pow(y[key], 2)
	}

	return dotProduct / (math.Sqrt(X) * math.Sqrt(Y))
}

// Euclidean compute the euclidean distance.
func Euclidean(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		//nolint: gomnd
		sum += math.Pow((x[key] - y[key]), 2)
	}

	return math.Sqrt(sum)
}

// Manhattan compute the manhatan distance.
func Manhattan(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		sum += math.Abs(x[key] - y[key])
	}

	return sum
}

// Camberra compute the camberra distance.
func Camberra(x, y map[string]float64) float64 {
	var sum float64

	for key := range x {
		sum += math.Abs(x[key]-y[key]) / math.Abs(x[key]+y[key])
	}

	return sum
}

// Chebyshev compute the chebyshev distance.
func Chebyshev(x, y map[string]float64) float64 {
	res := make([]float64, 0, len(x))

	for key := range x {
		res = append(res, math.Abs(x[key]-y[key]))
	}

	sort.Float64s(res)

	return res[len(res)-1]
}

// Minkowski compute the minkowski distance.
func Minkowski(x, y map[string]float64, p float64) float64 {
	var sum float64

	for key := range x {
		abs := math.Abs(x[key] - y[key])
		sum += math.Pow(abs, p)
	}

	return math.Pow(sum, 1/p)
}

func ComputeDistance(name string, x, y map[string]float64) float64 {
	switch name {
	case "cosine":
		return Cosine(x, y)
	case "manhattan":
		return Manhattan(x, y)
	case "canberra":
		return Camberra(x, y)
	case "chebyshev":
		return Chebyshev(x, y)
	case "minkowski":
		//nolint: gomnd
		return Minkowski(x, y, 6)
	default:
		return Euclidean(x, y)
	}
}

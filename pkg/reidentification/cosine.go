package reidentification

import "math"

func CosineSimilarity(x, y map[string]float64) float64 {
	var dotProduct, X, Y float64

	//nolint: gomnd
	for key := range x {
		dotProduct += x[key] * y[key]
		X += math.Pow(x[key], 2)
		Y += math.Pow(y[key], 2)
	}

	return dotProduct / (math.Sqrt(X) * math.Sqrt(Y))
}

// func CountVectorizer(x, y map[string]interface{}) (xVect map[interface{}]int, yVect map[interface{}]int) {
// 	for _, val := range x {
// 		xVect[val]++
// 		yVect[val] = 0
// 	}

// 	for _, val := range y {
// 		xVect[val] = 0
// 		yVect[val]++
// 	}
// }

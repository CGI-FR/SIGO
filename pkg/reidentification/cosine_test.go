package reidentification_test

import (
	"math"
	"testing"

	"github.com/cgi-fr/sigo/pkg/reidentification"
	"github.com/stretchr/testify/assert"
)

func TestCosineSimilarity(t *testing.T) {
	t.Parallel()

	X := map[string]float64{
		"a": 14, "b": 6, "c": 18, "d": 52.1,
		"e": 21, "f": 36.48, "g": 12, "h": 39,
	}
	Y := map[string]float64{
		"a": 14.17, "b": 6, "c": 4, "d": 12.86,
		"e": 54, "f": 49, "g": 7.2, "h": 27.12,
	}

	cosine := reidentification.CosineSimilarity(X, Y)

	assert.InDelta(t, 0.7622963735959798, cosine, math.Pow10(-15))
}

// func TestCountVectorizer(t *testing.T) {
// 	t.Parallel()

// 	X := map[string]interface{}{"x": 14, "y": 18, "z": "ville"}
// 	Y := map[string]interface{}{"x": 7, "y": 2, "z": "ville"}
// 	Z := map[string]interface{}{"x": 3, "y": 15, "z": "mer"}
// 	U := map[string]interface{}{"x": 16, "y": 2, "z": "mer"}
// 	V := map[string]interface{}{"x": 11, "y": 9, "z": "ville"}
// 	W := map[string]interface{}{"x": 7, "y": 18, "z": "campagne"}

// 	// cosine(X,Y)

// 	vectorX := map[interface{}]int{14: 1, 18: 1, "ville": 1, 7: 0, 2: 0}
// 	vectorY := map[interface{}]int{14: 0, 18: 0, "ville": 1, 7: 1, 2: 1}
// }

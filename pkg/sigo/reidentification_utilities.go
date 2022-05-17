// Copyright (C) 2022 CGI France
//
// This file is part of SIGO.
//
// SIGO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SIGO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SIGO.  If not, see <http://www.gnu.org/licenses/>.

package sigo

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"

	"github.com/cgi-fr/jsonline/pkg/cast"
)

// MapItoMapF convert a map[string]interface{} to a map[string]float64.
func MapItoMapF(m map[string]interface{}) map[string]float64 {
	mFloat := make(map[string]float64)

	for key, value := range m {
		var val float64
		switch t := value.(type) {
		case int:
			val = float64(t)
		case string:
			//nolint: gomnd
			val, _ = strconv.ParseFloat(t, 64)
		case float32:
			val = float64(t)
		case json.Number:
			val, _ = t.Float64()
		case float64:
			val = t
		}

		mFloat[key] = val
	}

	return mFloat
}

// CountValues returns a map with the number of occurrences for each sensitive data value.
func CountValues(sensitives []string) map[string]int {
	count := make(map[string]int)
	for _, val := range sensitives {
		count[val]++
	}

	return count
}

// IsUnique returns if the string slice contains unique values or not.
func IsUnique(sensitives []string) bool {
	count := CountValues(sensitives)

	return len(count) == 1
}

// Returns the list of values present in the cluster for each qi.
func ListValues(data []map[string]interface{}, s []string) (mapValues map[string][]interface{}) {
	mapValues = make(map[string][]interface{})

	for _, mapData := range data {
		for key, val := range mapData {
			for _, sens := range s {
				if key != sens {
					v, _ := cast.ToFloat64(val)
					mapValues[key] = append(mapValues[key], v)
				}
			}
		}
	}

	return mapValues
}

// ScaleData returns the all scaled data.
func ScaleData(data []map[string]interface{}, s []string) (scaledData []map[string]interface{}) {
	listValues := ListValues(data, s)

	for _, originalMap := range data {
		scaledMap := make(map[string]interface{})

		for key, val := range originalMap {
			for _, sens := range s {
				if key != sens {
					v, _ := cast.ToFloat64(val)
					scaledMap[key] = Scaling2(v, listValues[key])
				} else {
					// nolint: forcetypeassert
					scaledMap[key] = val.(string)
				}
			}
		}

		scaledData = append(scaledData, scaledMap)
	}

	return scaledData
}

// Scaling returns the scaled value to ensure the mean and the standard deviation to be 0 and 1, respectively.
func Scaling2(value interface{}, listValues []interface{}) float64 {
	listVals := SliceToFloat64(listValues)
	// Standardization
	return (value.(float64) - Mean(listVals)) / Std(listVals)
}

// SliceToFloat64 convert a slice of interface into a slice of float64.
func SliceToFloat64(slice []interface{}) (res []float64) {
	for _, elt := range slice {
		res = append(res, elt.(float64))
	}

	return res
}

// RoundFloat round a float with a given precision.
func RoundFloat(val float64, precision uint) float64 {
	//nolint: gomnd
	ratio := math.Pow(10, float64(precision))

	return math.Round(val*ratio) / ratio
}

func TopSimilarity(s map[float64]interface{}) (float64, interface{}) {
	scores := []float64{}

	for similarity := range s {
		scores = append(scores, similarity)
	}

	sort.Sort(sort.Reverse(sort.Float64Slice(scores)))

	// best score
	top := scores[0]

	return scores[0], s[top]
}

// Unique returns if the slice contains unique map[string]interface{} or not.
func Unique(slice []map[string]interface{}, qi []string) bool {
	tmp := make(map[string]int)

	for _, data := range slice {
		val := ""

		for i, q := range qi {
			s, _ := cast.ToString(data[q])
			if i == 0 {
				//nolint: forcetypeassert
				val += s.(string)
			} else {
				val += "-" + s.(string)
			}
		}

		tmp[val]++
	}

	return len(tmp) == 1
}

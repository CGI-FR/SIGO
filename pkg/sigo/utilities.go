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
func IsUnique(sensitives map[string][]interface{}) map[string]bool {
	unique := make(map[string]bool)

	for key, vals := range sensitives {
		count := CountValues(SliceString(vals))
		unique[key] = len(count) == 1
	}

	return unique
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

// Returns the list of values present in the slice of map[string]interface{}.
func ListValues(data []map[string]interface{}, s []string) (mapValues map[string][]interface{}) {
	mapValues = make(map[string][]interface{})

	for _, mapData := range data {
		for key, val := range mapData {
			if key != s[0] && key != s[1] {
				v, _ := cast.ToFloat64(val)
				mapValues[key] = append(mapValues[key], v)
			}
		}
	}

	return mapValues
}

// SliceToFloat64 convert a slice of interface into a slice of float64.
func SliceToFloat64(slice []interface{}) (res []float64) {
	for _, elt := range slice {
		val, _ := cast.ToFloat64(elt)
		res = append(res, val.(float64))
	}

	return res
}

// SliceString convert a slice of interface into a slice of string.
func SliceString(slice []interface{}) (res []string) {
	for _, elt := range slice {
		val, _ := cast.ToString(elt)
		res = append(res, val.(string))
	}

	return res
}

// Scale returns the scaled value to ensure the mean and the standard deviation to be 0 and 1, respectively.
func Scale(value interface{}, mean float64, std float64) float64 {
	val, _ := cast.ToFloat64(value)
	// Standardization
	return (val.(float64) - mean) / std
}

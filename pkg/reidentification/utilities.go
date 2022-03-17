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

package reidentification

import (
	"encoding/json"
	"strconv"
)

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

func RemoveDuplicate(floatSlice []float64) []float64 {
	keys := make(map[float64]bool)
	list := []float64{}

	for _, val := range floatSlice {
		if _, value := keys[val]; !value {
			keys[val] = true

			list = append(list, val)
		}
	}

	return list
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

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

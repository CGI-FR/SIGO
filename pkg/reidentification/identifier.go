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
	"fmt"
	"strconv"
	"strings"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Identifier struct {
	metric        string
	original      *[]map[string]interface{}
	masked        *[]map[string]interface{}
	groupedMasked *[]map[string]interface{}
}

func NewIdentifier(distance string) Identifier {
	return Identifier{
		metric:        distance,
		original:      &[]map[string]interface{}{},
		masked:        &[]map[string]interface{}{},
		groupedMasked: &[]map[string]interface{}{},
	}
}

type IdentifiedRecord struct {
	original  map[string]interface{}
	sensitive []string
}

// Record returns the original sigo.Record of IdentifiedRecord.
func (ir IdentifiedRecord) Record() sigo.Record {
	record := jsonline.NewRow()
	qi := []string{}

	for key, val := range ir.original {
		record.Set(key, val)
		qi = append(qi, key)
	}

	record.Set("sensitive", ir.sensitive)

	return infra.NewJSONLineRecord(&record, &qi, &[]string{"sensitive"})
}

// IsEmpty check if IdentifiedRecord is empty.
func (ir IdentifiedRecord) IsEmpty() bool {
	return ir.sensitive[0] == ""
}

// ReturnGroup returns anonymized data without duplicate individuals.
func (id Identifier) ReturnGroup() *[]map[string]interface{} {
	return id.groupedMasked
}

// SaveData saves data in Identifier according to the type of dataset (original, anonymized).
func (id Identifier) SaveData(dataset sigo.RecordSource, typeData string) {
	var sink *infra.SliceDictionariesSink

	switch typeData {
	case "original":
		sink = infra.NewSliceDictionariesSink(id.original)
	case "anonymized":
		sink = infra.NewSliceDictionariesSink(id.masked)
	}

	for dataset.Next() {
		err := sink.Collect(dataset.Value())
		if err != nil {
			fmt.Println("Cannot collect data")
		}
	}
}

// GroupMasked groups anonymized data of the same value
// {"x":7,"y":6.67,"z":"a"}, {"x":7,"y":6.67,"z":"a"}, {"x":7,"y":6.67,"z":"a"}
// returns {"x":7,"y":6.67,"z":"a"}
// {"x":3,"y":7,"z":"b"}, {"x":3,"y":7,"z":"a"}, {"x":3,"y":7,"z":"c"}
// returns {"x":3,"y":7,"z":""}.
func (id Identifier) GroupMasked(qi, s []string) {
	sink := infra.NewSliceDictionariesSink(id.groupedMasked)
	// map containing for each tuple of quasi-identifier the list of sensitive data
	tmp := make(map[string][]string)

	for _, record := range *id.masked {
		val := ""

		for i, q := range qi {
			if i == 0 {
				val += string(record[q].(json.Number))
			} else {
				val += "-" + string(record[q].(json.Number))
			}
		}

		for _, sens := range s {
			tmp[val] = append(tmp[val], record[sens].(string))
		}
	}

	for val, sensitives := range tmp {
		vals := strings.Split(val, "-")
		row := jsonline.NewRow()

		for i, q := range qi {
			//nolint: gomnd
			val, _ := strconv.ParseFloat(vals[i], 64)
			row.Set(q, val)
		}

		if IsUnique(sensitives) {
			row.Set(s[0], sensitives[0])
		} else {
			row.Set(s[0], "")
		}

		record := infra.NewJSONLineRecord(&row, &qi, &s)

		err := sink.Collect(record)
		if err != nil {
			fmt.Println("Cannot collect data")
		}
	}
}

// Identify returns an IdentifiedRecord if an anonymized record matches the original record.
func (id Identifier) Identify(scaledData map[string]interface{}, originalData map[string]interface{},
	qi, s []string) IdentifiedRecord {
	sims := NewSimilarities(id.metric)
	i := 0

	scaledAnonymized := ScaleData(*id.groupedMasked, s)

	// for each anonymized scaled data
	for _, record := range scaledAnonymized {
		sim := NewSimilarity(i, record, qi, s)
		X := MapItoMapF(scaledData)
		Y := MapItoMapF(sim.qi)
		// we calculate the distance with the original data
		score := ComputeDistance(id.metric, X, Y)
		sim.AddScore(score)

		sims.Add(sim)
		i++
	}

	// we collect the most similar data to the original data
	top := sims.TopSimilarity()

	// if the most similar data have sensitive data
	sensitive := top.sensitive

	return IdentifiedRecord{original: originalData, sensitive: sensitive}
}

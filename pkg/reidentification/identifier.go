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
	"strings"

	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Identifier struct {
	metric        string
	k             int
	masked        *[]map[string]interface{}
	groupedMasked *[]map[string]interface{}
}

func NewIdentifier(distance string, k int) Identifier {
	return Identifier{
		metric: distance, k: k, masked: &[]map[string]interface{}{},
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
	return len(ir.sensitive) == 0
}

// SaveMasked saves anonymized data in memory.
func (id Identifier) SaveMasked(maskedDataset sigo.RecordSource) {
	sink := infra.NewSliceDictionariesSink(id.masked)

	for maskedDataset.Next() {
		err := sink.Collect(maskedDataset.Value())
		if err != nil {
			fmt.Println("Cannot collect data")
		}
	}
}

// GroupMasked groups anonymized data of the same value.
func (id Identifier) GroupMasked(maskedDataset sigo.RecordSource, qi, s []string) {
	sink := infra.NewSliceDictionariesSink(id.groupedMasked)
	tmp := make(map[string][]string)

	for maskedDataset.Next() {
		record := maskedDataset.Value()
		val := ""

		for i, q := range qi {
			if i == 0 {
				val += string(record.Row()[q].(json.Number))
			} else {
				val += "-" + string(record.Row()[q].(json.Number))
			}
		}

		for _, sens := range s {
			tmp[val] = append(tmp[val], record.Row()[sens].(string))
		}
	}

	for val, sensitives := range tmp {
		vals := strings.Split(val, "-")
		row := jsonline.NewRow()

		for i, q := range qi {
			row.Set(q, vals[i])
		}

		if IsUnique(sensitives) {
			row.Set(s[0], sensitives[0])
		} else {
			row.Set(s[0], nil)
		}

		record := infra.NewJSONLineRecord(&row, &qi, &s)

		err := sink.Collect(record)
		if err != nil {
			fmt.Println("Cannot collect data")
		}
	}
}

func (id Identifier) ReturnsGroup() *[]map[string]interface{} {
	return id.groupedMasked
}

// Identify returns an IdentifiedRecord if an anonymized record matches the original record.
func (id Identifier) Identify(originalRec sigo.Record, maskedDataset sigo.RecordSource,
	qi, s []string) IdentifiedRecord {
	// map containing the original sigo.record
	x := make(map[string]interface{})

	for _, q := range qi {
		x[q] = originalRec.Row()[q]
	}

	sims := NewSimilarities(id.metric)
	i := 0

	// for each anonymized data
	for _, record := range *id.masked {
		sim := NewSimilarity(i, record, qi, s)
		X := MapItoMapF(x)
		Y := MapItoMapF(sim.qi)
		// we calculate the distance with the original data
		score := ComputeDistance(id.metric, X, Y)
		sim.AddScore(score)

		sims.Add(sim)
		i++
	}

	// we collect the k most similar data to the original data
	top := sims.TopSimilarity(id.k)
	// if the n most similar data have the same value of sensitive data
	sensitive, risk := Recover(top.Slice())

	// then we can reassociate the anonymized data with its sensitive data
	if risk {
		return IdentifiedRecord{original: x, sensitive: sensitive}
	}

	return IdentifiedRecord{original: x, sensitive: []string{}}
}

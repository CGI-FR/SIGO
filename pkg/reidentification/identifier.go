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
	"github.com/cgi-fr/jsonline/pkg/jsonline"
	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"
)

type Identifier struct {
	metric Distance
	k      int
}

func NewIdentifier(name string, k int) Identifier {
	var metric Distance

	switch name {
	case "cosine":
		metric = NewCosineSimilarity()
	case "manhattan":
		metric = NewManhattanDistance()
	case "canberra":
		metric = NewCanberraDistance()
	case "chebyshev":
		metric = NewChebyshevDistance()
	case "minkowski":
		//nolint: gomnd
		metric = NewMinkowskiDistance(3)
	default:
		metric = NewEuclideanDistance()
	}

	return Identifier{metric: metric, k: k}
}

type IdentifiedRecord struct {
	original  map[string]interface{}
	sensitive []string
}

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

func (ir IdentifiedRecord) IsEmpty() bool {
	return len(ir.sensitive) == 0
}

func (id Identifier) Identify(originalRec sigo.Record, maskedDataset sigo.RecordSource,
	qi, s []string) IdentifiedRecord {
	x := make(map[string]interface{})

	for _, q := range qi {
		x[q] = originalRec.Row()[q]
	}

	sims := NewSimilarities(id.metric)
	i := 0

	for maskedDataset.Next() {
		sim := NewSimilarity(i, maskedDataset.Value(), maskedDataset.QuasiIdentifer(), maskedDataset.Sensitive())

		X := MapItoMapF(x)
		Y := MapItoMapF(sim.qi)
		score := id.metric.Compute(X, Y)
		sim.AddScore(score)

		sims.Add(sim)
		i++
	}

	top := sims.TopSimilarity(id.k)
	sensitive, risk := Recover(top.Slice())

	if risk {
		return IdentifiedRecord{original: x, sensitive: sensitive}
	}

	return IdentifiedRecord{original: x, sensitive: []string{}}
}

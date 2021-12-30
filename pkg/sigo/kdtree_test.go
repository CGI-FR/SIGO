// Copyright (C) 2021 CGI France
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

package sigo_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	t.Parallel()

	kdtree := sigo.NewKDTreeFactory().New(1, 1)

	row := jsonline.NewRow()
	row.Set("x", 0)

	record := infra.NewJSONLineRecord(&row, &[]string{"x"})

	kdtree.Add(record)

	clusters := kdtree.Clusters()

	assert.Equal(t, 1, len(clusters))
	assert.Equal(t, 1, len(clusters[0].Records()))
	assert.Equal(t, 0, clusters[0].Records()[0].Row()["x"])
}

// nolint: funlen
func TestAddNRecords(t *testing.T) {
	t.Parallel()

	type test struct {
		k int
		n int
		d int
		s int64
	}

	tests := []test{
		{k: 1, n: 10, d: 1, s: 10},
		{k: 2, n: 10, d: 1, s: 10},
		{k: 3, n: 10, d: 1, s: 10},
		{k: 4, n: 10, d: 1, s: 10},
		{k: 5, n: 10, d: 1, s: 10},
		{k: 6, n: 10, d: 1, s: 10},
		{k: 1, n: 10, d: 1, s: 10},
		{k: 2, n: 20, d: 1, s: 10},
		{k: 3, n: 30, d: 1, s: 10},
		{k: 4, n: 40, d: 1, s: 10},
		{k: 5, n: 100, d: 1, s: 10},
		{k: 6, n: 1000, d: 1, s: 10},
		{k: 1, n: 10, d: 1, s: 10},
		{k: 2, n: 10, d: 2, s: 10},
		{k: 3, n: 10, d: 3, s: 10},
		{k: 4, n: 10, d: 4, s: 10},
		{k: 5, n: 10, d: 5, s: 10},
		{k: 6, n: 10, d: 6, s: 10},
		{k: 1, n: 10, d: 1, s: 10},
		{k: 2, n: 20, d: 2, s: 10},
		{k: 3, n: 30, d: 3, s: 10},
		{k: 4, n: 40, d: 4, s: 10},
		{k: 5, n: 100, d: 5, s: 10},
		{k: 6, n: 1000, d: 6, s: 10},
		{k: 1, n: 10, d: 1, s: 5},
		{k: 2, n: 10, d: 1, s: 5},
		{k: 3, n: 10, d: 1, s: 5},
		{k: 4, n: 10, d: 1, s: 5},
		{k: 5, n: 10, d: 1, s: 5},
		{k: 6, n: 10, d: 1, s: 5},
		{k: 1, n: 10, d: 1, s: 5},
		{k: 2, n: 20, d: 1, s: 5},
		{k: 3, n: 30, d: 1, s: 5},
		{k: 4, n: 40, d: 1, s: 5},
		{k: 5, n: 100, d: 1, s: 5},
		{k: 6, n: 1000, d: 1, s: 5},
		{k: 1, n: 10, d: 1, s: 5},
		{k: 2, n: 10, d: 2, s: 5},
		{k: 3, n: 10, d: 3, s: 5},
		{k: 4, n: 10, d: 4, s: 5},
		{k: 5, n: 10, d: 5, s: 5},
		{k: 6, n: 10, d: 6, s: 5},
		{k: 1, n: 10, d: 1, s: 5},
		{k: 2, n: 20, d: 2, s: 5},
		{k: 3, n: 30, d: 3, s: 5},
		{k: 4, n: 40, d: 4, s: 5},
		{k: 5, n: 100, d: 5, s: 5},
		{k: 6, n: 1000, d: 6, s: 5},
	}

	// nolint: paralleltest
	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			t.Parallel()

			k := tc.k
			N := tc.n
			D := tc.d

			rand.Seed(tc.s)

			kdtree := sigo.NewKDTreeFactory().New(k, 1)
			rows := []jsonline.Row{}

			for i := 0; i < N; i++ {
				// nolint: gosec
				x := rand.Intn(N)
				// nolint: gosec
				y := rand.Intn(N)
				for j := 0; j < D; j++ {
					row := jsonline.NewRow()
					row.Set("x", x)
					row.Set("y", y)
					rows = append(rows, row)
				}
			}

			// fmt.Printf("%v\n", rows)

			for i := 0; i < N; i++ {
				record := infra.NewJSONLineRecord(&rows[i], &[]string{"x"})

				kdtree.Add(record)
			}

			clusters := kdtree.Clusters()

			// fmt.Println(kdtree.(sigo.KDTree).String())

			for i := 0; i < len(clusters); i++ {
				assert.True(t, len(clusters[i].Records()) >= k)
			}
		})
	}
}

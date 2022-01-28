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

package sigo_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/cgi-fr/sigo/internal/infra"
	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	t.Parallel()

	kdtree := sigo.NewKDTreeFactory().New(1, 1, 1)

	row := jsonline.NewRow()
	row.Set("x", 0)

	record := infra.NewJSONLineRecord(&row, &[]string{"x"}, &[]string{})

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

			kdtree := sigo.NewKDTreeFactory().New(k, 1, 1)
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
				record := infra.NewJSONLineRecord(&rows[i], &[]string{"x"}, &[]string{})

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

func TestAddClusterInfos(t *testing.T) {
	t.Parallel()

	kdtree := sigo.NewKDTreeFactory().New(3, 1, 2)
	rows := []jsonline.Row{}
	x := []int{20, 10, 12, 24, 8, 16, 15, 27, 25, 11, 49, 2, 35, 34, 21}
	y := []int{10, 12, 4, 21, 38, 16, 26, 18, 30, 19, 21, 12, 14, 7, 5}
	z := []string{"a", "b", "a", "a", "c", "c", "b", "a", "b", "c", "a", "c", "a", "b", "a"}
	expectedPath := []string{
		"root-u-l", "root-l-l", "root-l-l", "root-u-u", "root-l-u",
		"root-l-u", "root-l-u", "root-u-u", "root-u-u", "root-l-u",
		"root-u-u", "root-l-l", "root-u-l", "root-u-l", "root-u-l",
	}

	for i := range x {
		row := jsonline.NewRow()
		row.Set("x", x[i])
		row.Set("y", y[i])
		row.Set("z", z[i])
		rows = append(rows, row)
	}

	for i := range x {
		record := infra.NewJSONLineRecord(&rows[i], &[]string{"x", "y"}, &[]string{})

		kdtree.Add(record)
	}

	kdtree.Build()
	clusters := kdtree.Clusters()

	for _, cluster := range clusters {
		for _, record := range cluster.Records() {
			for i := range rows {
				result, _ := rows[i].Export()
				if reflect.DeepEqual(result.(map[string]interface{}), record.Row()) {
					assert.Equal(t, cluster.ID(), expectedPath[i])
				}
			}
		}
	}
}

func BenchmarkClustering(b *testing.B) {
	N := []int{1000, 10000, 25000, 50000, 100000}

	for _, n := range N {
		b.Run(fmt.Sprintf("input_size_%d", n), func(b *testing.B) {
			rand.Seed(time.Now().UnixNano())

			kdtree := sigo.NewKDTreeFactory().New(3, 1, 1)
			rows := []jsonline.Row{}

			for i := 0; i < n; i++ {
				row := jsonline.NewRow()
				// nolint: gosec
				x := rand.Intn(n)
				row.Set("x", x)
				// nolint: gosec
				y := rand.Intn(n)
				row.Set("y", y)
				rows = append(rows, row)
			}

			for j := 0; j < n; j++ {
				record := infra.NewJSONLineRecord(&rows[j], &[]string{"x", "y"}, &[]string{})

				kdtree.Add(record)
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				kdtree.Build()
			}
		})
	}
}

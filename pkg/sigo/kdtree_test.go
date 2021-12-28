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

	"github.com/cgi-fr/sigo/pkg/sigo"

	"github.com/cgi-fr/jsonline/pkg/jsonline"

	"github.com/stretchr/testify/assert"
)

func TestAddRecord(t *testing.T) {
	t.Parallel()

	kdtree := sigo.NewKDTreeFactory().New(1, 1)

	row := jsonline.NewRow()
	row.Set("x", 0)

	record := NewJSONLineRecord(&row, &[]string{"x"})

	kdtree.Add(record)

	clusters := kdtree.Clusters()

	assert.Equal(t, 1, len(clusters))
	assert.Equal(t, 1, len(clusters[0].Records()))
	assert.Equal(t, 0, clusters[0].Records()[0].Row()["x"])
}

func TestAddNRecords(t *testing.T) {
	t.Parallel()

	k := 2
	N := 8

	rand.Seed(10)

	kdtree := sigo.NewKDTreeFactory().New(k, 1)
	rows := []jsonline.Row{}

	for i := 0; i < N; i++ {
		row := jsonline.NewRow()
		// nolint: gosec
		row.Set("x", rand.Float32())
		rows = append(rows, row)
	}

	fmt.Printf("%v\n", rows)

	for i := 0; i < N; i++ {
		record := NewJSONLineRecord(&rows[i], &[]string{"x"})

		kdtree.Add(record)
	}

	clusters := kdtree.Clusters()

	fmt.Println(kdtree.(sigo.KDTree).String())

	assert.Truef(t, len(clusters) >= N/k, "#clusters(%d) != %d", len(clusters), N/k)

	for i := 0; i < len(clusters); i++ {
		assert.True(t, len(clusters[i].Records()) >= k)
	}
}

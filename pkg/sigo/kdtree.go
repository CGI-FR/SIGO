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

package sigo

import (
	"fmt"
	"math"
	"sort"
	"strings"

	over "github.com/Trendyol/overlog"
)

func NewKDTreeFactory() KDTreeFactory {
	return KDTreeFactory{}
}

type KDTreeFactory struct{}

func (f KDTreeFactory) New(k int, l int, dim int) Generalizer {
	// nolint: exhaustivestruct
	tree := KDTree{k: k, l: l, dim: dim, clusterID: make(map[string]int)}
	root := newNode(&tree, "root", 0)
	root.validate()
	tree.root = &root

	return tree
}

type KDTree struct {
	k         int
	l         int
	root      *node
	dim       int
	clusterID map[string]int
}

func (t KDTree) Add(r Record) {
	t.root.add(r)
}

func (t KDTree) Build() {
	t.root.build()
}

func (t KDTree) Clusters() []Cluster {
	return t.root.clusters()
}

func (t KDTree) String() string {
	return t.root.string(0)
}

func newNode(tree *KDTree, path string, rot int) node {
	return node{
		tree:        tree,
		cluster:     []Record{},
		clusterPath: path,
		subNodes:    []node{},
		pivot:       []float32{},
		valid:       false,
		rot:         rot % tree.dim,
	}
}

type node struct {
	tree        *KDTree
	cluster     []Record
	clusterPath string
	subNodes    []node
	pivot       []float32
	valid       bool
	rot         int
}

func (n *node) add(r Record) {
	n.cluster = append(n.cluster, r)
}

func (n *node) build() {
	if n.isValid() && len(n.cluster) >= 2*n.tree.k {
		// rollback to simple node
		lower, upper, valide := n.split()
		if !valide {
			return
		}

		// log.Info().Msgf("new pivot: %v", n.pivot)
		// log.Info().Str("node", lower.string(0)).Msg("new node")
		// log.Info().Str("node", upper.string(0)).Msg("new node")

		lower.validate()
		upper.validate()

		n.subNodes = []node{
			lower,
			upper,
		}

		n.cluster = nil

		n.subNodes[0].build()
		n.subNodes[1].build()
	}
}

func (n *node) split() (node, node, bool) {
	sort.SliceStable(n.cluster, func(i int, j int) bool {
		return n.cluster[i].QuasiIdentifer()[n.rot] < n.cluster[j].QuasiIdentifer()[n.rot]
	})

	n.pivot = nil
	lower := newNode(n.tree, n.clusterPath+"-l", n.rot+1)
	upper := newNode(n.tree, n.clusterPath+"-u", n.rot+1)
	lowerSize := 0
	upperSize := 0
	previous := n.cluster[0]

	for _, row := range n.cluster {
		if lowerSize < len(n.cluster)/2 || row.QuasiIdentifer()[n.rot] == previous.QuasiIdentifer()[n.rot] {
			lower.add(row)
			previous = row
			lowerSize++
		} else {
			if n.pivot == nil {
				n.pivot = row.QuasiIdentifer()
			}
			upper.add(row)
			upperSize++
		}
	}

	return lower, upper, upperSize >= n.tree.k && lower.wellLDiv() && upper.wellLDiv()
}

func (n *node) Records() []Record {
	if n.cluster != nil {
		return n.cluster
	}

	return []Record{}
}

func (n *node) ID() string {
	return n.clusterPath
}

func (n *node) clusters() []Cluster {
	if n.cluster != nil {
		return []Cluster{n}
	}

	return append(n.subNodes[0].clusters(), n.subNodes[1].clusters()...)
}

func (n *node) string(offset int) string {
	if n.cluster != nil {
		result := "["
		for _, rec := range n.cluster {
			// result += fmt.Sprintf("%v ", rec.QuasiIdentifer()[n.rot])
			result += fmt.Sprintf("%v ", rec.QuasiIdentifer())
		}

		result += "]"

		return result
	}

	return fmt.Sprintf("{\n%s pivot: %v,\n%s rot: %v, \n%s n0: %s,\n%s n1: %s,\n%s}",
		strings.Repeat(" ", offset),
		n.pivot[n.rot],
		strings.Repeat(" ", offset),
		n.rot,
		strings.Repeat(" ", offset),
		n.subNodes[0].string(offset+1),
		strings.Repeat(" ", offset),
		n.subNodes[1].string(offset+1),
		strings.Repeat(" ", offset),
	)
}

func (n *node) validate() {
	n.valid = true
}

func (n *node) isValid() bool {
	return n.valid
}

func (n node) wellLDiv() bool {
	var f func([]Record, int) float64

	if b, _ := over.MDC().Get("entropy"); b.(bool) {
		f = entropy
	} else {
		f = logQ
	}

	rec := n.cluster[0]
	for i := 0; i < len(rec.Sensitives()); i++ {
		e := f(n.cluster, i)
		if e < math.Log(float64(n.tree.l)) {
			return false
		}
	}

	return true
}

func entropy(clus []Record, sens int) float64 {
	frequency := make(map[interface{}]int)

	for _, rec := range clus {
		val := rec.Sensitives()[sens]
		frequency[val]++
	}

	var e float64
	for _, val := range frequency {
		e += (float64(val) / float64(len(clus))) * math.Log(float64(val)/float64(len(clus)))
	}

	return e
}

func logQ(clus []Record, sens int) float64 {
	frequency := make(map[interface{}]int)

	for _, rec := range clus {
		val := rec.Sensitives()[sens]
		frequency[val] = 1
	}

	var e float64
	for _, val := range frequency {
		e += float64(val)
	}

	return e
}

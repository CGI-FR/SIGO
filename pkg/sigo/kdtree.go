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
	"sort"
	"strings"
)

func NewKDTreeFactory() KDTreeFactory {
	return KDTreeFactory{}
}

type KDTreeFactory struct{}

func (f KDTreeFactory) New(k int, l int) Generalizer {
	// nolint: exhaustivestruct
	tree := KDTree{k: k, l: l}
	root := newNode(&tree)
	root.validate()
	tree.root = &root

	return tree
}

type KDTree struct {
	k    int
	l    int
	root *node
}

func (t KDTree) Add(r Record) {
	t.root.add(r)
}

func (t KDTree) Clusters() []Cluster {
	return t.root.clusters()
}

func (t KDTree) String() string {
	return t.root.string(0)
}

func newNode(tree *KDTree) node {
	// nolint: exhaustivestruct
	return node{tree: tree, valid: false, cluster: []Record{}}
}

type node struct {
	tree     *KDTree
	cluster  []Record
	subNodes []node
	pivot    []float32
	valid    bool
}

func (n *node) add(r Record) {
	if n.cluster == nil {
		if r.QuasiIdentifer()[0] < n.pivot[0] {
			n.subNodes[0].add(r)
		} else {
			n.subNodes[1].add(r)
		}

		return
	}

	n.cluster = append(n.cluster, r)

	if n.isValid() && len(n.cluster) >= 2*n.tree.k {
		// rollback to simple node
		lower, upper, shouldReturn := n.newMethod()
		if shouldReturn {
			return
		}

		lower.validate()
		upper.validate()

		n.subNodes = []node{
			lower,
			upper,
		}
		n.cluster = nil
	}
}

func (n *node) newMethod() (node, node, bool) {
	sort.SliceStable(n.cluster, func(i int, j int) bool {
		return n.cluster[i].QuasiIdentifer()[0] < n.cluster[j].QuasiIdentifer()[0]
	})

	n.pivot = nil
	lower := newNode(n.tree)
	upper := newNode(n.tree)
	lowerSize := 0
	upperSize := 0
	previous := n.cluster[0]

	for _, row := range n.cluster {
		if lowerSize < len(n.cluster)/2 || row.QuasiIdentifer()[0] == previous.QuasiIdentifer()[0] {
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

	return lower, upper, upperSize < n.tree.k
}

func (n *node) Records() []Record {
	if n.cluster != nil {
		return n.cluster
	}

	return []Record{}
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
			result += fmt.Sprintf("%v ", rec.QuasiIdentifer()[0])
		}

		result += "]"

		return result
	}

	return fmt.Sprintf("{\n%s pivot: %v,\n%s n0: %s,\n%s n1: %s,\n%s}",
		strings.Repeat(" ", offset),
		n.pivot[0],
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

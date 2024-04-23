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

package sigo

import (
	"fmt"
	"math"
	"sort"
	"strings"

	over "github.com/Trendyol/overlog"
	"github.com/rs/zerolog/log"
)

func NewKDTreeFactory() KDTreeFactory {
	return KDTreeFactory{}
}

type KDTreeFactory struct{}

func (f KDTreeFactory) New(k int, l int, dim int, qi []string) Generalizer {
	// nolint: exhaustivestruct
	tree := KDTree{k: k, l: l, dim: dim, clusterID: make(map[string]int), qi: qi}
	root := NewNode(&tree, "root", 0)
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
	qi        []string
}

func NewKDTree(k, l, dim int, clusterID map[string]int) KDTree {
	//nolint: exhaustivestruct
	return KDTree{k: k, l: l, dim: dim, clusterID: clusterID}
}

// Add add a record to the tree (root node).
func (t KDTree) Add(r Record) {
	t.root.Add(r)
}

// Build starts building the tree.
func (t KDTree) Build() error {
	return t.root.build()
}

// Clusters returns the list of clusters in the tree.
func (t KDTree) Clusters() []Cluster {
	return t.root.Clusters()
}

// String returns the tree in a literary way.
func (t KDTree) String() string {
	return t.root.string(0)
}

// nolint: revive, golint
func NewNode(tree *KDTree, path string, rot int) node {
	return node{
		tree:        tree,
		cluster:     []Record{},
		clusterPath: path,
		subNodes:    []node{},
		pivot:       []float64{},
		valid:       false,
		rot:         rot % tree.dim,
	}
}

type node struct {
	tree        *KDTree
	cluster     []Record
	clusterPath string
	subNodes    []node
	pivot       []float64
	valid       bool
	rot         int
}

// Add add a record to the node.
func (n *node) Add(r Record) {
	n.cluster = append(n.cluster, r)
}

// incRot increments the value of rot which refers to the dimension on which we build our nodes.
func (n *node) incRot() {
	n.rot = (n.rot + 1) % n.tree.dim
}

// build creates nodes.
func (n *node) build() error {
	log.Debug().
		Str("Dimension", n.tree.qi[n.rot]).
		Str("Path", n.clusterPath).
		Int("Size", len(n.cluster)).
		Msg("Cluster:")

	if n.isValid() && len(n.cluster) >= 2*n.tree.k {
		// rollback to simple node
		var (
			lower, upper node
			valide       bool
			err          error
		)

		for i := 1; i <= n.tree.dim; i++ {
			lower, upper, valide, err = n.split()
			if err != nil {
				return err
			} else if !valide {
				n.incRot()
			} else {
				break
			}
		}

		if !valide {
			return nil
		}

		lower.validate()
		upper.validate()

		n.subNodes = []node{
			lower,
			upper,
		}

		n.cluster = nil
		err = n.subNodes[0].build()
		if err != nil {
			return err
		}
		err = n.subNodes[1].build()
		if err != nil {
			return err
		}
	}

	return nil
}

// split creates 2 subnodes by ordering the node and splitting in order to have 2 equal parts
// and all elements having the same value in the same subnode.
func (n *node) split() (node, node, bool, error) {
	var globalError error

	less := func(i, j int) bool {
		valueI, err := n.cluster[i].QuasiIdentifer()
		if err != nil {
			// Stocker l'erreur dans la variable globale
			globalError = err
			return false
		}
		valueJ, err := n.cluster[j].QuasiIdentifer()
		if err != nil {
			globalError = err
			return false
		}
		return valueI[n.rot] < valueJ[n.rot]
	}

	sort.SliceStable(n.cluster, less)
	if globalError != nil {
		return node{}, node{}, false, globalError
	}

	n.pivot = nil
	lower := NewNode(n.tree, n.clusterPath+"-l", n.rot+1)
	upper := NewNode(n.tree, n.clusterPath+"-u", n.rot+1)

	lowerSize := 0
	upperSize := 0
	previous := n.cluster[0]

	for _, row := range n.cluster {
		rowValue, _ := row.QuasiIdentifer()
		previousValue, _ := previous.QuasiIdentifer()
		// equal subnodes and all elements having the same value in the same subnode
		if lowerSize < len(n.cluster)/2 || rowValue[n.rot] == previousValue[n.rot] {
			lower.Add(row)
			previous = row
			lowerSize++
		} else {
			if n.pivot == nil {
				n.pivot = rowValue
			}
			upper.Add(row)
			upperSize++
		}
	}

	return lower, upper, upperSize >= n.tree.k && lower.wellLDiv() && upper.wellLDiv(), nil
}

// Records returns the list of records in the node.
func (n *node) Records() []Record {
	if n.cluster != nil {
		return n.cluster
	}

	return []Record{}
}

// ID return the path of the node.
func (n *node) ID() string {
	return n.clusterPath
}

// Clusters returns the list of clusters in the node.
func (n *node) Clusters() []Cluster {
	if n.cluster != nil {
		return []Cluster{n}
	}

	return append(n.subNodes[0].Clusters(), n.subNodes[1].Clusters()...)
}

// string retunrs the node information (pivot, dimension, subnodes).
func (n *node) string(offset int) string {
	if n.cluster != nil {
		result := "["
		for _, rec := range n.cluster {
			// result += fmt.Sprintf("%v ", rec.QuasiIdentifer()[n.rot])
			recValue, _ := rec.QuasiIdentifer()
			result += fmt.Sprintf("%v ", recValue)
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

// validate validates the node.
func (n *node) validate() {
	n.valid = true
}

// isValid returns if the node is valid or not.
func (n *node) isValid() bool {
	return n.valid
}

// weelDiv returns if the node respects the l-diversity (https://en.wikipedia.org/wiki/L-diversity).
// (https://tel.archives-ouvertes.fr/tel-01783967/document p.29/30).
func (n node) wellLDiv() bool {
	var f func([]Record, int) float64

	var condition float64

	if b, ok := over.MDC().Get("entropy"); !ok || !b.(bool) {
		// Distinct l-diversity
		f = logQ
		condition = float64(n.tree.l)
	} else {
		// Entropy l-diversity
		f = entropy
		condition = math.Log(float64(n.tree.l))
	}

	rec := n.cluster[0]
	for i := 0; i < len(rec.Sensitives()); i++ {
		e := f(n.cluster, i)
		if e < condition {
			return false
		}
	}

	return true
}

// entropy returns the value of the entropy for the cluster clus and the sensible attribute sens.
// (https://tel.archives-ouvertes.fr/tel-01783967/document p.30).
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

// logQ returns the number of represented values in the cluster clus for the sensible attribute sens.
// (https://tel.archives-ouvertes.fr/tel-01783967/document p.29).
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

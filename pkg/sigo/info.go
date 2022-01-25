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

func NewSequenceDebugger(key string) SequenceDebugger { return SequenceDebugger{map[string]int{}, key} }

type SequenceDebugger struct {
	cache map[string]int
	key   string
}

type InfosRecord struct {
	original Record
	infos    map[string]interface{}
}

func (ir InfosRecord) QuasiIdentifer() []float32 {
	return ir.original.QuasiIdentifer()
}

func (ir InfosRecord) Sensitives() []interface{} {
	return ir.original.Sensitives()
}

func (ir InfosRecord) Row() map[string]interface{} {
	original := ir.original.Row()
	for k, v := range ir.infos {
		original[k] = v
	}

	return original
}

func (d SequenceDebugger) id(c Cluster) int {
	count := len(d.cache)
	if d.cache[c.ID()] == 0 {
		d.cache[c.ID()] = count + 1
	}

	return d.cache[c.ID()]
}

func (d SequenceDebugger) Information(rec Record, cluster Cluster) Record {
	infos := make(map[string]interface{})

	infos[d.key] = d.id(cluster)

	return InfosRecord{original: rec, infos: infos}
}

func NewNoDebugger() Debugger { return NoDebugger{} }

type NoDebugger struct{}

func (d NoDebugger) Information(rec Record, cluster Cluster) Record {
	return rec
}

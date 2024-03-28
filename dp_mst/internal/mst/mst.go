package mst

import (
	cmn "dp_mst/internal/common"
	"sort"
)

// MF-Sets Functions ------------------------

func father(i int64, id []int64) int64 {
	for i != id[i] {
		id[i] = id[id[i]]
		i = id[i]
	}
	return i
}

func Father(i int64, id []int64) int64 {
	return father(i, id)
}

func unite(p, q int64, id []int64) {
	i := father(p, id)
	j := father(q, id)
	id[i] = j
}

func Unite(p, q int64, id []int64) {
	unite(p, q, id)
}

// Minimum Spanning Tree Functions ------------------------

func Kruskal(root map[int64]*cmn.Graph, mst cmn.Graph) cmn.Graph {
	sizeG := 0
	for _, adje := range root {
		sizeG += len(*adje)
	}
	localG := make(cmn.Graph, 0, sizeG)
	for _, adje := range root {
		localG = append(localG, *adje...)
	}

	sort.Slice(localG, func(i, j int) bool {
		return localG[i].W < localG[j].W || (localG[i].W == localG[j].W && localG[i].X < localG[j].X)
	})

	m := make(map[int64]int64)
	var id []int64
	var cc int64
	cc = 0
	returnMST := make(cmn.Graph, 0, len(mst)+len(localG))

	idMST := 0
	idLocalG := 0
	for idMST < len(mst) || idLocalG < len(localG) {
		var e cmn.Edge
		if idMST < len(mst) && idLocalG < len(localG) {
			if mst[idMST].W < localG[idLocalG].W {
				e = mst[idMST]
				idMST++
			} else {
				e = localG[idLocalG]
				idLocalG++
			}
		} else if idMST < len(mst) {
			e = mst[idMST]
			idMST++
		} else {
			e = localG[idLocalG]
			idLocalG++
		}

		// Check wheter each edge is already in graph
		_, ok1 := m[e.X]
		_, ok2 := m[e.Y]

		if !ok1 || !ok2 { // If one of the vertices is not, add the edge
			returnMST = append(returnMST, e)
			if !ok1 {
				m[e.X] = cc
				id = append(id, cc)
				cc++
			}
			if !ok2 {
				m[e.Y] = cc
				id = append(id, cc)
				cc++
			}

			unite(m[e.X], m[e.Y], id)
			m[e.X] = id[m[e.X]]
			m[e.Y] = id[m[e.Y]]
		} else if father(m[e.X], id) != father(m[e.Y], id) {
			returnMST = append(returnMST, e)
			unite(m[e.X], m[e.Y], id)
			m[e.X] = id[m[e.X]]
			m[e.Y] = id[m[e.Y]]
		}
	}

	return returnMST
}

// Frederickson 85 ------------------------
type GraphP struct {
}

// Implements the data structure for dynamic trees
// proposed by Sleator and Trajan.
// DOI: https://doi.org/10.1016/0022-0000(83)90006-5
type NodeST struct {
	external, reversed              bool
	parent, left, right, head, tail *NodeST
	netmin, netcost                 float64
}

func (n *NodeST) Path() *NodeST {
	var m *NodeST
	m = n
	for m.parent != nil {
		m = m.parent
	}
	return m
}

func (n *NodeST) Head() *NodeST {
	if n.reversed {
		return n.tail
	} else {
		return n.head
	}
}

func (n *NodeST) Tail() *NodeST {
	if n.reversed {
		return n.head
	} else {
		return n.tail
	}
}

func (n *NodeST) Before() *NodeST {
	var m *NodeST
	m = n
	for (!m.reversed && m.parent.left == m) ||
		(m.reversed && m.parent.right == m) ||
		m.parent == nil {
		m = m.parent
	}
	if m.parent == nil {
		panic("Asked before the head")
	}
	m = m.parent
	if m.parent.left != nil {
		m = m.left
		for (!m.reversed && m.right != nil) ||
			(m.reversed && m.left != nil) {
			if !m.reversed {
				m = m.right
			} else {
				m = m.left
			}
		}
	}
	return m
}

func (n *NodeST) After() *NodeST {
	var m *NodeST
	m = n
	for (!m.reversed && m.parent.right == m) ||
		(m.reversed && m.parent.left == m) ||
		m.parent == nil {
		m = m.parent
	}
	if m.parent == nil {
		panic("Asked before the head")
	}
	m = m.parent
	if m.parent.right != nil {
		m = m.right
		for (!m.reversed && m.left != nil) ||
			(m.reversed && m.right != nil) {
			if !m.reversed {
				m = m.left
			} else {
				m = m.right
			}
		}
	}
	return m
}

func (n *NodeST) Cost() float64 {
	return 0
}

func (n *NodeST) MinCost() float64 {
	return 0
}

func (n *NodeST) Update(x float64) {
	n.netmin += x
}

func (n *NodeST) Reverse() {
	n.reversed = !n.reversed
}

// Assumes g is a MST
func Frederickson(g *cmn.Graph) {

}

func csearch() {

}

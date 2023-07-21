package mst

import (
	"sort"
	cmn "dp_mst/internal/common"
)

// MF-Sets Functions ------------------------

func father(i int64, id []int64) int64 {
	for i != id[i] {
		id[i] = id[id[i]]
		i = id[i]
	}
	return i
}

func unite(p,q int64, id []int64) {
	i := father(p, id)
	j := father(q, id)
	id[i] = j
}

// Minimum Spanning Tree Functions ------------------------

func Kruskal(root map[int64]*cmn.Graph, mst cmn.Graph) cmn.Graph {
	for _, adje := range root {
		mst = append(mst, *adje...)
	}

	sort.Slice(mst, func(i, j int) bool {
		return mst[i].W < mst[j].W || (mst[i].W == mst[j].W && mst[i].X < mst[j].X)
	})

	m := make(map[int64]int64)
	var id []int64
	var cc int64
	cc = 0
	var returnMST cmn.Graph

	for _, e := range mst {
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

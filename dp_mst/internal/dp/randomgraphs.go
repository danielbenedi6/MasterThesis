package dp

import (
	cmn "dp_mst/internal/common"
	"math/rand"
)

func ErdosRenyi(probMST float64, probEdge float64, N int64, out out_comm, Fsize int) {
	var i int64 = 0
	var j int64 = i + 1
	var r cmn.Request
	
	for i < N {
		if probMST > rand.Float64() {
			r.Op = cmn.Operation("kmst")
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
			out.Req <- r 
			empty := make(cmn.Graph, 0)
			out.Graph <- empty


			r.Op = cmn.Operation("graph")
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
			out.Req <- r 
			empty = make(cmn.Graph, 0)
			out.Graph <- empty
		} else {
			if probEdge < rand.Float64() {
				r.Op = cmn.Operation("insert")
				r.E = cmn.Edge{X: i, Y: j, W: rand.Float64()}
				out.Req <- r 
			}

			j++
			if j == N {
				i++
				j = i + 1
			}
		}
	}
}
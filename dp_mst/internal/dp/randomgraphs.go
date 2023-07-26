package dp

import (
	cmn "dp_mst/internal/common"
	"math/rand"
)

func ErdosRenyi(probMST float64, probEdge float64, N int64, seed int64, out out_comm, Fsize int) {
	var i int64 = 0
	var j int64 = i + 1
	var r cmn.Request

	random := rand.New(rand.NewSource(seed))

	for i < N {
		if probMST > random.Float64() {
			r.Op = cmn.KMST
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
			out.Req <- r
			empty := make(cmn.Graph, 0)
			out.Graph <- empty

			r.Op = cmn.GraphOp
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
			out.Req <- r
			empty = make(cmn.Graph, 0)
			out.Graph <- empty
		} else {
			if probEdge < random.Float64() {
				r.Op = cmn.Insert
				r.E = cmn.Edge{X: i, Y: j, W: random.Float64()}
				out.Req <- r
			}

			j++
			if j == N {
				i++
				j = i + 1
			}
		}
	}

	r.Op = cmn.EOF
	r.E = cmn.Edge{X: -1, Y: -1, W: 0}
	out.Req <- r
	empty := make(cmn.Graph, 0)
	out.Graph <- empty
}

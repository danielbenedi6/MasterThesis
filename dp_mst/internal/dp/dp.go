package dp

import (
	"bufio"
	cmn "dp_mst/internal/common"
	"dp_mst/internal/mst"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const channelSize = 5000

type in_comm struct {
	Req   <-chan cmn.Request
	Graph <-chan cmn.Graph
}

type out_comm struct {
	Req   chan<- cmn.Request
	Graph chan<- cmn.Graph
}

func input(istream string, out out_comm, Fsize int) {
	file, err := os.Open(istream + ".requests")
	cmn.CheckError(err)
	defer file.Close()
	cmn.CheckError(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var r cmn.Request
	for {
		scanner.Scan()
		op := scanner.Text()
		err = scanner.Err()
		cmn.CheckError(err)
		if err == io.EOF || op == "" {
			break
		}

		op_int, err := strconv.ParseInt(op, 10, 64)
		if err != nil {
			break
		}
		cmn.CheckError(err)

		r.Op = cmn.Operation(op_int)

		if r.Op == cmn.KMST || r.Op == cmn.GraphOp || r.Op == cmn.EOF {
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
		} else {
			scanner.Scan()
			node1 := scanner.Text()
			scanner.Scan()
			node2 := scanner.Text()
			scanner.Scan()
			weight := scanner.Text()
			err = scanner.Err()
			cmn.CheckError(err)

			r.E.X, err = strconv.ParseInt(node1, 10, 32)
			r.E.Y, err = strconv.ParseInt(node2, 10, 32)
			r.E.W, err = strconv.ParseFloat(weight, 64)
		}

		r.Normalize()

		out.Req <- r

		if r.Op == cmn.EOF {
			break
		}

		if r.Op == cmn.KMST || r.Op == cmn.GraphOp {
			empty := make(cmn.Graph, 0)
			out.Graph <- empty
		}

	}
	close(out.Req)
	close(out.Graph)
}

func output(istream string, in in_comm, end chan<- struct{}) {
	for {
		r, ok := <-in.Req
		if !ok {
			break
		}

		switch r.Op {
		case cmn.GraphOp:
			g, _ := <-in.Graph
			fmt.Println("Graph", g)
		case cmn.KMST:
			g, _ := <-in.Graph
			fmt.Println("MST", g)
		case cmn.EOF:
			break
		default: //something's wrong
			fmt.Println("Unknown operation in output")
			break
		}
	}
	end <- struct{}{}
}

func generator(in in_comm, out out_comm, Fsize int) {
	filter_count := 0
	for {
		r, ok := <-in.Req
		if !ok {
			break
		}
		switch r.Op {
		case cmn.Insert, cmn.Update:
			out_req := make(chan cmn.Request, channelSize)
			out_grph := make(chan cmn.Graph, channelSize)
			new_out := out_comm{Req: out_req, Graph: out_grph}
			new_in := in_comm{Req: out_req, Graph: out_grph}

			go filter(filter_count, in, new_out, r.E, Fsize)
			filter_count++
			in = new_in
		case cmn.Delete:
			// Do nothing, asked to delete unexistant edge
		case cmn.GraphOp, cmn.KMST:
			g, _ := <-in.Graph

			out.Req <- r
			out.Graph <- g
		case cmn.EOF:
			out.Req <- r
			break
		default: //something's wrong
			fmt.Println("Unknown operation in generator")
			break
		}
	}
	close(out.Req)
	close(out.Graph)
}

func filter(id int, in in_comm, out out_comm, e cmn.Edge, Fsize int) {
	s := map[int64]struct{}{}

	int_comm := make(chan cmn.Request, 4096)

	go filter_worker(id, in_comm{int_comm, in.Graph}, out.Graph)

	for {
		r, ok := <-in.Req
		if !ok {
			break
		}

		switch r.Op {
		case cmn.Insert, cmn.Update:
			if _, ok = s[r.E.X]; ok {
				int_comm <- r
			} else if len(s) < Fsize {
				s[r.E.X] = struct{}{}
				int_comm <- r
			} else {
				out.Req <- r
			}
		case cmn.Delete:
			if _, ok = s[r.E.X]; ok {
				int_comm <- r
			} else {
				out.Req <- r
			}
		case cmn.KMST:
			int_comm <- r
		case cmn.GraphOp:
			int_comm <- r
		case cmn.EOF:
			out.Req <- r
			break
		}
	}

	close(out.Req)
	close(out.Graph)
}

func filter_worker(id int, in in_comm, out chan<- cmn.Graph) {

	// Initialize memory
	root := make(map[int64]*cmn.Graph)

	for {
		r, ok := <-in.Req
		if !ok {
			break
		}

		switch r.Op {
		case cmn.Insert, cmn.Update:
			if _, ok = root[r.E.X]; ok {
				root[r.E.X].InsertUpdate(r.E)
			} else {
				root[r.E.X] = &cmn.Graph{r.E}
			}
		case cmn.Delete:
			if _, ok = root[r.E.X]; ok {
				root[r.E.X].Delete(r.E)
			}
		case cmn.KMST:
			g, _ := <-in.Graph
			g = mst.Kruskal(root, g)
			out <- g
		case cmn.GraphOp:
			g, _ := <-in.Graph
			for _, adje := range root {
				g = append(g, *adje...)
			}

			local_root := make(map[int64]cmn.Graph)
			for id, adje := range root {
				local_root[id] = *adje
			}

			out <- g
		case cmn.EOF:
			break
		}
	}

	close(out)
}

func Start(istream string, Fsize int) {

	file_req := make(chan cmn.Request, channelSize)
	file_grph := make(chan cmn.Graph, channelSize)
	file_gen := out_comm{Req: file_req, Graph: file_grph}
	gen := in_comm{Req: file_req, Graph: file_grph}

	gen_req := make(chan cmn.Request, channelSize)
	gen_grph := make(chan cmn.Graph, channelSize)
	gen_out := out_comm{Req: gen_req, Graph: gen_grph}
	out := in_comm{Req: gen_req, Graph: gen_grph}

	end := make(chan struct{})

	start := time.Now()
	go input(istream, file_gen, Fsize)
	go generator(gen, gen_out, Fsize)
	go output(istream, out, end)
	<-end
	t := time.Since(start)
	fmt.Println("TotalExecutionTime,", t, ",", t.Microseconds(), "μs,", t.Milliseconds(), "ms ,", t.Seconds(), "s")
}

func RandomStart(Fsize int, Nmin, Nmax, Ndelta, Ngraph, Nrep int64, Dmin, Dmax, Ddelta float64) {
	f, _ := os.Create("RandomStart.csv")
	defer f.Close()
	f.WriteString("N;D;seed;microseconds\n")

	for N := Nmin; N < Nmax; N += Ndelta {
		for D := Dmin; D < Dmax; D += Ddelta {
			for i := Ngraph; i > 0; i-- {
				var seed int64 = rand.Int63()
				for j := Nrep; j > 0; j-- {
					input_req := make(chan cmn.Request, channelSize)
					input_grph := make(chan cmn.Graph, channelSize)
					input_chan := out_comm{Req: input_req, Graph: input_grph}
					data_chan := in_comm{Req: input_req, Graph: input_grph}

					gen_req := make(chan cmn.Request, channelSize)
					gen_grph := make(chan cmn.Graph, channelSize)
					gen_out := out_comm{Req: gen_req, Graph: gen_grph}
					out := in_comm{Req: gen_req, Graph: gen_grph}

					end := make(chan struct{})

					start := time.Now()
					go ErdosRenyi(0.05, D, N, seed, input_chan, Fsize)
					go generator(data_chan, gen_out, Fsize)
					go output("", out, end)
					<-end
					t := time.Since(start)

					fmt.Fprintf(f, "%d;%.4f;%d;%d\n", N, D, seed, t.Microseconds())
				}
			}
		}
	}
}

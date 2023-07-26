package main

import (
	"bufio"
	cmn "dp_mst/internal/common"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const file_extension = ".requests"
const extension_len = len(file_extension)
const repetitions = 100

type in_comm struct {
	Req   <-chan cmn.Request
	Graph <-chan cmn.Graph
}

type out_comm struct {
	Req   chan<- cmn.Request
	Graph chan<- cmn.Graph
}

func input_fmt(istream string, out out_comm, Fsize int) {
	file, err := os.Open(istream + file_extension)
	cmn.CheckError(err)
	defer file.Close()
	cmn.CheckError(err)

	var r cmn.Request
	for {
		_, err = fmt.Fscanf(file, "%s", &r.Op)
		if err == io.EOF {
			break
		}
		cmn.CheckError(err)

		if r.Op == cmn.KMST || r.Op == cmn.GraphOp || r.Op == cmn.EOF {
			r.E = cmn.Edge{X: -1, Y: -1, W: 0}
		} else {
			_, err = fmt.Fscanf(file, "%d %d %f\n", &r.E.X, &r.E.Y, &r.E.W)
			cmn.CheckError(err)
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

func split(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' || data[i] == '\n' {
			return i + 1, data[:i], nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func input_bufio(istream string, out out_comm, Fsize int) {
	file, err := os.Open(istream + file_extension)
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
			<-in.Graph
			//fmt.Println("Graph", g)
		case cmn.KMST:
			<-in.Graph
			//fmt.Println("MST", g)
		case cmn.EOF:
			break
		}
	}
	end <- struct{}{}
}

func main() {
	var dir string

	flag.StringVar(&dir, "dir", "", "Specify root directory of test files")
	flag.Parse()

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(dir + ".csv")
	defer f.Close()
	f.WriteString("N;reader;microseconds\n")

	for _, entry := range entries {
		if strings.Contains(entry.Name(), file_extension) {
			size := entry.Name()[:len(entry.Name())-extension_len]
			fmt.Println(dir + "/" + entry.Name())
			for i := 0; i < repetitions; i++ {
				file_req := make(chan cmn.Request, 5000)
				file_grph := make(chan cmn.Graph, 5000)
				file_gen := out_comm{Req: file_req, Graph: file_grph}
				out := in_comm{Req: file_req, Graph: file_grph}

				end := make(chan struct{})

				start := time.Now()
				go input_fmt(dir+"/"+size, file_gen, 0)
				go output(size, out, end)
				<-end
				t := time.Since(start)

				f.WriteString(size + ";fmt;" + strconv.FormatInt(t.Microseconds(), 10) + "\n")

				file_req = make(chan cmn.Request, 5000)
				file_grph = make(chan cmn.Graph, 5000)
				file_gen = out_comm{Req: file_req, Graph: file_grph}
				out = in_comm{Req: file_req, Graph: file_grph}

				end = make(chan struct{})

				start = time.Now()
				go input_bufio(dir+"/"+size, file_gen, 0)
				go output(size, out, end)
				<-end
				t = time.Since(start)

				f.WriteString(size + ";bufio;" + strconv.FormatInt(t.Microseconds(), 10) + "\n")
			}
		}
	}
}

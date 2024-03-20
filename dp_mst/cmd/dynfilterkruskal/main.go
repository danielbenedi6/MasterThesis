package main

import (
	"bufio"
	cmn "dp_mst/internal/common"
	mst "dp_mst/internal/mst"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type solverMST func([]cmn.Edge, *[]cmn.Edge, *[]int64, map[int64]int64, *int64)

func Kruskal(G []cmn.Edge, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
	sort.Slice(G, func(i, j int) bool {
		return G[i].W < G[j].W
	})

	for _, e := range G {
		// Check wheter each cmn.Edge is already in graph
		_, ok1 := M[e.X]
		_, ok2 := M[e.Y]

		if !ok1 || !ok2 { // If one of the vertices is not, add the cmn.Edge
			*MST = append(*MST, e)
			if !ok1 {
				M[e.X] = *CC
				*P = append(*P, *CC)
				*CC++
			}
			if !ok2 {
				M[e.Y] = *CC
				*P = append(*P, *CC)
				*CC++
			}

			mst.Unite(M[e.X], M[e.Y], *P)
			M[e.X] = (*P)[M[e.X]]
			M[e.Y] = (*P)[M[e.Y]]
		} else if mst.Father(M[e.X], *P) != mst.Father(M[e.Y], *P) {
			*MST = append(*MST, e)
			mst.Unite(M[e.X], M[e.Y], *P)
			M[e.X] = (*P)[M[e.X]]
			M[e.Y] = (*P)[M[e.Y]]
		}
	}
}

func filterKruskal(G []cmn.Edge, left, right int64, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
	size := right - left + 1
	if size < int64(len(G)*5/100) || size < 1000 {
		Kruskal(G[left:right], MST, P, M, CC)
		return
	}

	p := G[rand.Int63n(size)]

	l := left
	r := right
	// Do partition
	for l < r {
		for l < r && G[l].W <= p.W {
			l++
		}
		for l < r && G[r].W > p.W {
			r--
		}

		if l < r {
			G[l], G[r] = G[r], G[l]
		}
	}
	filterKruskal(G, left, r, MST, P, M, CC)

	// Do partition, but now from l and checking UnionFind ids
	r = right
	left = l + 1
	for l < r {
		//Condition to keep: mst.mst.Father(e.X, P) != mst.mst.Father(e.Y, P)
		for l < r && mst.Father(M[G[l].X], *P) != mst.Father(M[G[l].Y], *P) {
			l++
		}
		for l < r && mst.Father(M[G[r].X], *P) == mst.Father(M[G[r].Y], *P) {
			r--
		}

		if l < r {
			G[l], G[r] = G[r], G[l]
		}
	}

	if left < r {
		filterKruskal(G, left, r, MST, P, M, CC)
	}
}

func FilterKruskal(G []cmn.Edge, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
	filterKruskal(G, 0, int64(len(G)-1), MST, P, M, CC)
}

func ReadDynGraph(istream string, fMST solverMST) {
	file, err := os.Open(istream + ".requests")
	cmn.CheckError(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var G cmn.Graph
	for {
		scanner.Scan()
		op := scanner.Text()
		err = scanner.Err()
		cmn.CheckError(err)
		if err == io.EOF || op == "" {
			break
		}
		op_int, err := strconv.ParseInt(op, 10, 64)

		scanner.Scan()
		node1 := scanner.Text()
		scanner.Scan()
		node2 := scanner.Text()
		scanner.Scan()
		weight := scanner.Text()
		err = scanner.Err()
		cmn.CheckError(err)

		var E cmn.Edge
		E.X, err = strconv.ParseInt(node1, 10, 32)
		E.Y, err = strconv.ParseInt(node2, 10, 32)
		E.W, err = strconv.ParseFloat(weight, 64)

		if cmn.Operation(op_int) == cmn.Insert {
			G.InsertUpdate(E)
		} else if cmn.Operation(op_int) == cmn.Delete {
			G.Delete(E)
		}

		MST := make([]cmn.Edge, 0)
		P := make([]int64, 0)
		M := make(map[int64]int64)
		CC := int64(0)
		fMST(G, &MST, &P, M, &CC)
	}
}

func main() {
	var reps = flag.Int("repetitions", 1, "Number of times MST is computed for reproductibility")

	flag.Parse()

	filePath := "../DynGraphRepo/stats.csv"

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	var wg sync.WaitGroup

	csvReader := csv.NewReader(f)
	csvReader.Comma = ';'
	line, err := csvReader.Read()
	for err == nil {
		instance := line[0]

		times_fk := make(chan float64, *reps+1)
		times_kr := make(chan float64, *reps+1)
		for exp := 0; exp < *reps; exp++ {
			go func() {
				wg.Add(1)
				defer wg.Done()

				start := time.Now()
				ReadDynGraph("../DynGraphRepo/"+instance, FilterKruskal)
				elapsed := time.Since(start)
				times_fk <- float64(elapsed)
			}()
			go func() {
				wg.Add(1)
				defer wg.Done()

				start := time.Now()
				ReadDynGraph("../DynGraphRepo/"+instance, Kruskal)
				elapsed := time.Since(start)
				times_fk <- float64(elapsed)
			}()
		}
		wg.Wait()

		fmt.Printf("FilterKruskal,%s", instance)
		for exp := 0; exp < *reps; exp++ {
			time := <-times_fk
			fmt.Printf(",%.3f", time/1000.0)
		}
		fmt.Printf("\n")

		fmt.Printf("Kruskal,%s", instance)
		for exp := 0; exp < *reps; exp++ {
			time := <-times_kr
			fmt.Printf(",%.3f", time/1000.0)
		}
		fmt.Printf("\n")

		line, err = csvReader.Read()
	}
}

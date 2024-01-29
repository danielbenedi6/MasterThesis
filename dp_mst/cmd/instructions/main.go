package main

import (
	cmn "dp_mst/internal/common"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Edge = struct{ X, Y int }

func GenerateGraph(N int, p float64, seed int64) (map[Edge]float64, map[Edge]struct{}) {
	rand_gen := rand.New(rand.NewSource(seed))

	var G map[Edge]float64

	// Algorithm based on: https://doi.org/10.1103/PhysRevE.71.036113
	v := 1
	w := -1
	for v < N {
		r := rand_gen.Float64()
		w = w + 1 + int(math.Log(1-r)/math.Log(1-p))
		for w >= v && v < N {
			w = w - v
			v = v + 1
		}

		if v < N {
			// Add edge
			G[Edge{X: v, Y: w}] = rand_gen.Float64()
		}
	}

	var notG map[Edge]struct{}

	for i := 1; i < N; i++ {
		for j := i + 1; j < N; j++ {
			if _, ok := G[Edge{i, j}]; !ok {
				notG[Edge{i, j}] = struct{}{}
			}
		}
	}

	return G, notG
}

func RandEdge[K comparable, V float64 | struct{}](G map[K]V, rand_gen *rand.Rand) K {
	k := rand.Intn(len(G))
	var res K
	for edge := range G {
		if k == 0 {
			return edge
		}
		res = edge
		k--
	}
	return res
}

func Experiment(G_orig map[Edge]float64, notG_orig map[Edge]struct{}, N int, p float64, Fsize, id, exp, ops, rep int) int64 {

	// Perform deep copy to avoid modifying the references
	G := map[Edge]float64{}
	notG := map[Edge]struct{}{}
	for k, v := range G_orig {
		G[k] = v
	}
	for k, v := range notG_orig {
		notG[k] = v
	}

	seed := time.Now().UnixNano()
	rand_gen := rand.New(rand.NewSource(seed))
	dir := fmt.Sprintf("operations/%d_%.2f_%d_%d", N, p, Fsize, id)
	os.MkdirAll(dir, os.ModePerm)

	f, err := os.OpenFile(dir+"/"+strconv.Itoa(exp)+".requeests", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(f, "%d\n%d\n", cmn.LoadState, cmn.CurrTime)

	for i := 0; i < rep; i++ {
		// Insert operations
		for j := 0; j < ops; j++ {
			edge := RandEdge[Edge, struct{}](notG, rand_gen)
			delete(G, edge)
			w := rand_gen.Float64()
			G[edge] = w

			fmt.Fprintf(f, "%d %d %d %0.5f\n", cmn.Insert, edge.X, edge.Y, w)
		}
		// Delete operations
		for j := 0; j < ops; j++ {
			edge := RandEdge[Edge, float64](G, rand_gen)
			delete(G, edge)
			notG[edge] = struct{}{}
			fmt.Fprintf(f, "%d %d %d\n", cmn.Delete, edge.X, edge.Y)
		}
		fmt.Fprintf(f, "%d\n", cmn.KMST)
	}

	return seed
}

func main() {
	var operations = flag.Int("operations", 1, "Number of insert/delete operations")
	var reps = flag.Int("repetitions", 1, "Number of time insert/delete operations are done")
	var num_experiments = flag.Int("experiments", 1, "Number of experiments to generate")

	flag.Parse()

	filePath := "./create_stats.csv"

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	var wg sync.WaitGroup

	csvReader := csv.NewReader(f)
	line, err := csvReader.Read()
	for err == nil {
		N, _ := strconv.ParseInt(line[0], 10, 32)
		p, _ := strconv.ParseFloat(line[1], 64)
		Fsize, _ := strconv.ParseInt(line[2], 10, 32)
		id, _ := strconv.ParseInt(line[3], 10, 32)
		seed, _ := strconv.ParseInt(line[5], 10, 64)

		wg.Add(1)

		go func() {
			defer wg.Done()
			G, notG := GenerateGraph(int(N), p, seed)
			for exp := 0; exp < *num_experiments; exp++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					Experiment(G, notG, int(N), p, int(Fsize), int(id), exp, *operations, *reps)
				}()
			}
		}()

		line, err = csvReader.Read()
	}

	wg.Wait()
}

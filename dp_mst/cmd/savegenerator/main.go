package main

import (
	cmn "dp_mst/internal/common"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func GenerateGraph(N int, p float64, filters int, id int) int64 {
	seed := time.Now().UnixNano()
	rand_gen := rand.New(rand.NewSource(seed))
	dir := fmt.Sprintf("saves/%d_%.2f_%d_%d", N, p, filters, id)
	os.MkdirAll(dir, os.ModePerm)

	max_filters := N / int(filters)

	filter_struct := make([]map[int64]struct{}, 1, max_filters)
	filter_struct[0] = map[int64]struct{}{}
	worker_struct := make([]map[int64]*cmn.Graph, 1, max_filters)
	worker_struct[0] = map[int64]*cmn.Graph{}
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
			_, ok := filter_struct[len(filter_struct)-1][int64(v)]
			if !ok && len(filter_struct[len(filter_struct)-1]) < filters {
				filter_struct[len(filter_struct)-1][int64(v)] = struct{}{}
				worker_struct[len(worker_struct)-1][int64(v)] = &cmn.Graph{}
			} else if !ok {
				filter_struct = append(filter_struct, map[int64]struct{}{int64(v): {}})
				worker_struct = append(worker_struct, map[int64]*cmn.Graph{int64(v): {}})
			}

			worker_struct[len(worker_struct)-1][int64(v)].InsertUpdate(cmn.Edge{
				X: int64(v),
				Y: int64(w),
				W: rand_gen.Float64(),
			})
		}
	}

	// Once all the structs have been filled with the edges of the random graph
	// It is as simply as perform a savestate
	for id, s := range filter_struct {
		file, err := os.Create(dir + "/filter_" + strconv.Itoa(id) + ".bin")
		if err != nil {
			log.Fatal("Cannot open file for savestate of filer " + strconv.Itoa(id))
		}

		// Create a new gob encoder and use it to encode the person struct
		enc := gob.NewEncoder(file)
		if err := enc.Encode(s); err != nil {
			fmt.Println("Error encoding struct:", err)
			os.Exit(3)
		}

		file.Close()
	}

	for id, root := range worker_struct {
		file, err := os.Create(dir + "/worker_filter_" + strconv.Itoa(id) + ".bin")
		if err != nil {
			log.Fatal("Cannot open file for savestate of worker " + strconv.Itoa(id))
		}

		// Create a new gob encoder and use it to encode the person struct
		enc := gob.NewEncoder(file)
		if err := enc.Encode(root); err != nil {
			fmt.Println("Error encoding state of worker:", err)
			os.Exit(3)
		}

		file.Close()
	}

	return seed
}

func main() {
	var number_nodes = flag.Int("numbernodes", 100, "Number of nodes in the graph.")
	var edge_probability = flag.Float64("edgeprob", 0.2, "Probability of an edge to be added.")
	var nodes_filter = flag.Int("filters", 1, "Number of nodes per filter.")
	var repetitions = flag.Uint("repetitions", 5, "Number of graph repetitions.")

	flag.Parse()

	if *edge_probability < 0 || *edge_probability > 1 {
		fmt.Fprintf(os.Stderr, "Probability of an edge should be between 0 and 1. Value provided: %f\n", *edge_probability)
		os.Exit(2)
	}
	if *number_nodes < 2 {
		fmt.Fprintf(os.Stderr, "A graph needs at least 2 nodes. Value provided: %d\n", *number_nodes)
		os.Exit(2)
	}
	if *nodes_filter < 0 {
		fmt.Fprintf(os.Stderr, "The number of nodes per filter must be greater than 0. Value provided: %d\n", *nodes_filter)
		os.Exit(2)
	}

	f, err := os.OpenFile("./create_stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	for id := 0; id < int(*repetitions); id++ {
		t_ini := time.Now()
		seed := GenerateGraph(*number_nodes, *edge_probability, *nodes_filter, id)
		t_end := time.Now()

		fmt.Fprintf(f, "%d,%.2f,%d,%d,%s,%d\n", *number_nodes, *edge_probability, *nodes_filter, id, t_end.Sub(t_ini), seed)
	}

}

package main

import (
	"dp_mst/cmd/filterkruskal/parallel"
	cmn "dp_mst/internal/common"
	mst "dp_mst/internal/mst"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
	"time"
	"unsafe"
)

func GenerateGraph(N int64, p float64, seed int64) []cmn.Edge {
	rand_gen := rand.New(rand.NewSource(seed))

	var G []cmn.Edge

	// Algorithm based on: https://doi.org/10.1103/PhysRevE.71.036113
	v := int64(1)
	w := int64(-1)
	for v < N {
		r := rand_gen.Float64()
		w = w + 1 + int64(math.Log(1-r)/math.Log(1-p))
		for w >= v && v < N {
			w = w - v
			v = v + 1
		}

		if v < N {
			// Add cmn.Edge
			G = append(G, cmn.Edge{X: v, Y: w, W: rand_gen.Float64()})
		}
	}

	return G
}

func Kruskal(G []cmn.Edge, N int64, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
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

func FilterKruskal(G []cmn.Edge, N, left, right int64, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
	size := right - left + 1
	if size < N {
		Kruskal(G[left:right], N, MST, P, M, CC)
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
	FilterKruskal(G, N, left, r, MST, P, M, CC)

	// Do partition, but now from l and checking UnionFind ids
	r = right
	left = l + 1
	for l < r {
		//Condition to keep: mst.mst.Father(e.X, P) != mst.mst.Father(e.Y, P)
		for l < r && mst.Father(G[l].X, *P) != mst.Father(G[l].Y, *P) {
			l++
		}
		for l < r && mst.Father(G[r].X, *P) == mst.Father(G[r].Y, *P) {
			r--
		}

		if l < r {
			G[l], G[r] = G[r], G[l]
		}
	}

	if left < r {
		FilterKruskal(G, N, left, r, MST, P, M, CC)
	}
}

// func ParallelFilterKruskal(G []cmn.Edge, N, left, right int64, MST *[]cmn.Edge, P *[]int64, M map[int64]int64, CC *int64) {
// 	size := right - left + 1
// 	if size < N {
// 		Kruskal(G[left:right], N, MST, P, M, CC)
// 		return
// 	}
// 	p := G[rand.Int63n(size)]

// 	l := left
// 	r := right
// 	// Do partition
// 	for l < r {
// 		for l < r && G[l].W <= p.W {
// 			l++
// 		}
// 		for l < r && G[r].W > p.W {
// 			r--
// 		}

// 		if l < r {
// 			G[l], G[r] = G[r], G[l]
// 		}
// 	}
// 	ParallelFilterKruskal(G, N, left, r, MST, P, M, CC)

// 	numCPU := runtime.NumCPU() / bits.Len(uint(N*N))

// 	keep := make(chan cmn.Edge, right-l+1)

// 	var wg sync.WaitGroup
// 	wg.Add(numCPU)
// 	for grID := 0; grID < numCPU; grID++ {
// 		go func(grID int) {
// 			defer wg.Done()
// 			for i := int64(grID) + l + 1; i < right; i += int64(numCPU) {
// 				if mst.Father(G[i].X, *P) != mst.Father(G[i].Y, *P) {
// 					keep <- G[i]
// 				}
// 			}
// 		}(grID)
// 	}
// 	wg.Wait()

// 	left = l + 1
// 	r = left
// 	end := false
// 	for !end {
// 		select {
// 		case e, _ := <-keep:
// 			G[r] = e
// 			r++
// 		default:
// 			end = true
// 		}
// 	}

// 	if left < r {
// 		ParallelFilterKruskal(G, N, left, r, MST, P, M, CC)
// 	}
// }

func main() {
	N := int64(10)
	G := GenerateGraph(N, 0.75, 1321412333121)

	fmt.Println("G generated")

	acum1 := 0.0
	acum2 := 0.0
	reps := 1

	f, err := os.Create("parallel_profile.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for rep := 0; rep < reps; rep++ {
		start := time.Now()
		MST := make([]cmn.Edge, 0, N-1)
		P := make([]int64, N)
		for i := int64(0); i < N; i++ {
			P[i] = i
		}
		M := make(map[int64]int64)
		CC := int64(0)
		FilterKruskal(G, N, 0, int64(len(G)-1), &MST, &P, M, &CC)
		elapsed := time.Since(start)
		fmt.Println("Time: ", elapsed, " -> ", float64(elapsed))
		for idx, e := range MST {
			fmt.Println(idx, ": (", e.X, ":", e.Y, ") [", e.W, "]")
		}

		acum1 += float64(elapsed)

		start = time.Now()
		MST = make([]cmn.Edge, 0, N-1)
		P = make([]int64, N)
		for i := int64(0); i < N; i++ {
			P[i] = i
		}
		M = make(map[int64]int64)
		CC = int64(0)
		MSTptr := parallel.ParallelFilterKruskal(&G[0], int64(len(G)), N)
		MST = (*(*[10]cmn.Edge)(unsafe.Pointer(&MSTptr)))[: N-1 : N-1]
		elapsed = time.Since(start)
		fmt.Println("Time: ", elapsed, " -> ", float64(elapsed))

		for idx, e := range MST {
			fmt.Println(idx, ": (", e.X, ":", e.Y, ") [", e.W, "]")
		}

		acum2 += float64(elapsed)
	}

	acum1 /= float64(reps)
	fmt.Println("Avg. Filtered-Kruskal: ", acum1)
	acum2 /= float64(reps)
	fmt.Println("Avg. Parallel Filtered-Kruskal: ", acum2)
}

// func main() {
// 	var reps = flag.Int("repetitions", 1, "Number of times MST is computed for reproductibility")

// 	flag.Parse()

// 	filePath := "./create_stats.csv"

// 	f, err := os.Open(filePath)
// 	if err != nil {
// 		log.Fatal("Unable to read input file "+filePath, err)
// 	}
// 	defer f.Close()

// 	var wg sync.WaitGroup

// 	csvReader := csv.NewReader(f)
// 	line, err := csvReader.Read()
// 	for err == nil {
// 		N, _ := strconv.ParseInt(line[0], 10, 32)
// 		p, _ := strconv.ParseFloat(line[1], 64)
// 		//Fsize, _ := strconv.ParseInt(line[2], 10, 32)
// 		//id, _ := strconv.ParseInt(line[3], 10, 32)
// 		seed, _ := strconv.ParseInt(line[5], 10, 64)

// 		G := GenerateGraph(N, p, seed)
// 		times_fk := make(chan float64, *reps+1)
// 		times_kr := make(chan float64, *reps+1)
// 		for exp := 0; exp < *reps; exp++ {
// 			go func() {
// 				wg.Add(1)
// 				defer wg.Done()

// 				start := time.Now()
// 				MST := make([]cmn.Edge, 0)
// 				P := make([]int64, N)
// 				for i := int64(0); i < N; i++ {
// 					P[i] = i
// 				}
// 				M := make(map[int64]int64)
// 				CC := int64(0)
// 				FilterKruskal(G, N, 0, int64(len(G)-1), &MST, &P, M, &CC)
// 				elapsed := time.Since(start)
// 				times_fk <- float64(elapsed)
// 			}()
// 			go func() {
// 				wg.Add(1)
// 				defer wg.Done()

// 				start := time.Now()
// 				MST := make([]cmn.Edge, 0)
// 				P := make([]int64, N)
// 				for i := int64(0); i < N; i++ {
// 					P[i] = i
// 				}
// 				M := make(map[int64]int64)
// 				CC := int64(0)
// 				Kruskal(G, N, &MST, &P, M, &CC)
// 				elapsed := time.Since(start)
// 				times_kr <- float64(elapsed)
// 			}()
// 		}
// 		wg.Wait()

// 		fmt.Printf("FilterKruskal,%d,%.2f,", N, p)
// 		for exp := 0; exp < *reps; exp++ {
// 			time := <-times_fk
// 			fmt.Printf("%.3f, ", time/1000.0)
// 		}
// 		fmt.Printf("\n")

// 		fmt.Printf("Kruskal,%d,%.2f,", N, p)
// 		for exp := 0; exp < *reps; exp++ {
// 			time := <-times_kr
// 			fmt.Printf("%.3f, ", time/1000.0)
// 		}
// 		fmt.Printf("\n")

// 		line, err = csvReader.Read()
// 	}
// }

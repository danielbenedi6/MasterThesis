package main

import (
	"fmt"
	"math/rand"
	"os"

	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

// Based on bucket-test
// https://www.johndcook.com/Beautiful_Testing_ch10.pdf
func main() {
	var N int

	seed := time.Now().UnixNano()
	rand_gen := rand.New(rand.NewSource(seed))

	fmt.Scanf("%d", &N)

	k := int(N/5) - 1

	Δb := 1.0 / float64(k)

	f, err := os.Create("random_numbers.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	C := make([]int, k)
	//V := make([]float64, N)
	for i := 0; i < N; i++ {
		elem := rand_gen.Float64()

		C[int(elem/Δb)]++
		//V[i] = elem
		//fmt.Fprintln(f, V[i])
		fmt.Fprintln(f, elem)
	}
	//sort.Float64s(V)

	//maxKminus := 0.0
	//maxKplus := 0.0
	//for i := 0; i < N; i++ {
	//	if maxKplus < (float64(i)/float64(N) - V[i]) {
	//		maxKplus = float64(i)/float64(N) - V[i]
	//	}
	//
	//	if maxKminus < (V[i] - float64(i-1)/float64(N)) {
	//		maxKminus = V[i] - float64(i-1)/float64(N)
	//	}
	//}
	//maxKminus *= math.Sqrt(float64(N))
	//maxKplus *= math.Sqrt(float64(N))
	//
	//fmt.Println("K⁺ = ", maxKplus)
	//fmt.Println("K⁻ = ", maxKminus)

	χ2 := 0.0
	for j := 0; j < k; j++ {
		χ2 += (float64(C[j]) - float64(N)*Δb) * (float64(C[j]) - float64(N)*Δb) / (float64(N) * Δb)
	}

	dist := distuv.ChiSquared{K: float64(k - 1)}
	fmt.Printf("χ2_%d = %f\n", k-1, χ2)
	fmt.Printf("P(χ2_%d ≤ %f) = %f\n", k-1, χ2, dist.Prob(χ2))
	fmt.Printf("P(χ2_%d ≤ %f) = %f\n", k-1, χ2, dist.CDF(χ2))
}

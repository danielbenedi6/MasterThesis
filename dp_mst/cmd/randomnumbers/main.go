package main

import (
	"fmt"
	"math/rand"
	//"slices"
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

	k := int(N / 5)

	Δb := 1.0 / float64(k)

	C := make([]int, k)
	//V := make([]float64, N)
	for i := 0; i < N; i++ {
		elem := rand_gen.Float64()

		C[int(elem/Δb)]++
		//V[i] = elem
	}
	//slices.Sort(V)

	χ2 := 0.0
	for j := 0; j < k; j++ {
		χ2 += (float64(C[j]) - float64(N)*Δb) * (float64(C[j]) - float64(N)*Δb) / (float64(N) * Δb)
	}

	dist := distuv.ChiSquared{K: float64(k - 1)}
	fmt.Printf("χ2_%d = %f\n", k-1, χ2)
	fmt.Printf("P(χ2_%d ≤ %f) = %f\n", k-1, χ2, dist.Prob(χ2))
	fmt.Printf("P(χ2_%d ≤ %f) = %f\n", k-1, χ2, dist.CDF(χ2))
}

package main

import (
	dp "dp_mst/internal/dp"
	"flag"
	"fmt"
	"runtime"
)

func main() {
	//-----------------------
	//
	var istream string
	//Pprof variables -----------------------
	//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.StringVar(&istream, "file", "", "Specify input file. Default is emptyfile")
	flag.Parse()

	//---------------------------------------
	// Filter size (number of "root nodes" stored by each filter)
	var Fsize int
	fmt.Scan(&Fsize)

	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	fmt.Println("maxProcs: ", maxProcs, " numCPU: ", numCPU)

	dp.Start(istream, Fsize)
}

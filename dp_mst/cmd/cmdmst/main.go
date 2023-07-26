package main

import (
	dp "dp_mst/internal/dp"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	//-----------------------
	//
	var istream string
	//Pprof variables -----------------------
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.StringVar(&istream, "file", "", "Specify input file. Default is emptyfile")
	flag.Parse()

	//---------------------------------------
	// Filter size (number of "root nodes" stored by each filter)
	var Fsize int
	fmt.Scan(&Fsize)

	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	fmt.Println("maxProcs: ", maxProcs, " numCPU: ", numCPU)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	dp.Start(istream, Fsize)

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}

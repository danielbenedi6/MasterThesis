package common

import (
	"log"
)

type Operation string

const (
	Insert  Operation = "insert"
	Update            = "update"
	Delete            = "delete"
	KMST              = "kmst"
	EOF               = "eof"
	GraphOp           = "graph"
)

type Edge struct {
	X, Y int64
	W    float64
}

type Graph []Edge

type Request struct {
	Op Operation
	E  Edge
}

func (r *Request) Normalize() {
	if r.E.Y < r.E.X {
		r.E.X, r.E.Y = r.E.Y, r.E.X
	}
}

func (g *Graph) Delete(e Edge) {
	// For eficiency porpuse, instead of deleting
	// the edge, we substitute it with the last
	// and shrink it. If the order matters, another
	// data structure will be more suitable
	for i, edge := range *g {
		if edge.X == e.X && edge.Y == e.Y {
			(*g)[i] = (*g)[len(*g)-1]
			(*g) = (*g)[:len(*g)-1]
			return
		}
	}
}

func (g *Graph) InsertUpdate(e Edge) {
	for _, edge := range *g {
		if edge.X == e.X && edge.Y == e.Y {
			edge.W = e.W
			return
		}
	}
	(*g) = append(*g, e)
}

// Generic Functions -----------------------------

func CheckError(e error) {
	if e != nil {
		ChangeLogPrefiX()
		log.Fatalf("Fatal error --- %s\n", e.Error())
	}
}

func ChangeLogPrefiX() {
	// Set microseconds and full PATH of source code in logs
	log.SetFlags(log.Lmicroseconds | log.Llongfile)
}

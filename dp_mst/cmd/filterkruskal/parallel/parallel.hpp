// Golang Edge struct:
//   type Edge = struct {
//   	X, Y int64
//   	W    float64
//   }
typedef struct Edge
{
	long int X, Y;
	double W;
} Edge;

Edge* ParallelFilterKruskal(Edge* G, long int M, long int N);

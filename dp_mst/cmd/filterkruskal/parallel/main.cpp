#include <iostream>
#include <vector>
#include "parallel.hpp"

int main() {
    std::vector<Edge> graph{
        Edge{7,1,0.0241721},
        Edge{2,0,0.330454},
        Edge{2,1,0.259096},
        Edge{7,3,0.0758428},
        Edge{1,0,0.0878156},
        Edge{9,7,0.192849},
        Edge{7,6,0.228438},
        Edge{7,2,0.228641},
        Edge{3,2,0.230666},
        Edge{4,0,0.277872},
        Edge{6,4,0.300598},
        Edge{5,3,0.338712},
        Edge{7,4,0.426581},
        Edge{9,6,0.466282},
        Edge{8,7,0.433301},
        Edge{5,0,0.435342},
        Edge{9,0,0.595036},
        Edge{9,4,0.522043},
        Edge{6,1,0.554342},
        Edge{8,1,0.489894},
        Edge{8,2,0.593493},
        Edge{8,3,0.7457},
        Edge{8,4,0.644372},
        Edge{8,5,0.714453},
        Edge{6,5,0.783818},
        Edge{3,0,0.589017},
        Edge{9,8,0.827111},
        Edge{7,5,0.81317},
        Edge{8,6,0.975413},
        Edge{7,0,0.81147},
        Edge{4,1,0.940379}
    };
    int N = 10;

    Edge *MST = ParallelFilterKruskal(graph.data(), graph.size(), N);

    for(int i = 0; i < N-1; i++)
        std::cout << MST[i].X << " " << MST[i].Y << " " << MST[i].W << std::endl;

    delete[] MST;
}
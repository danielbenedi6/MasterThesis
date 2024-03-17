#include "parallel.hpp"
#include <unordered_map>
#include <vector>
#include <iostream>
void ParallelFK(
                    std::vector<Edge> G,
                    long int N,
                    std::vector<Edge>& MST,
                    std::vector<long int>& P, 
                    std::unordered_map<long int, long int>& M,
                    long int& CC
                );

void Kruskal(
                    std::vector<Edge> G,
                    long int N,
                    std::vector<Edge>& MST,
                    std::vector<long int>& P, 
                    std::unordered_map<long int, long int>& M,
                    long int& CC
                );

#include <random>
#include <algorithm>
#include <execution>

Edge* ParallelFilterKruskal(Edge* G, long int M, long int N) {
    std::vector<Edge> _G(G, G+M);
    //for(int i = 0; i < M; i++) {
    //    std::cout << "C++ G(" << _G[i].X << "," << _G[i].Y << ") [" << _G[i].W << "]" << std::endl;
    //}
    std::vector<Edge> _MST;
    std::vector<long int> ID;
    std::unordered_map<long int, long int> map;
    long int initCC = 0;

    _MST.reserve(N);
    ParallelFK(_G, N, _MST, ID, map, initCC);
	
    Edge *MST = new Edge[N-1];
    std::copy(_MST.begin(), _MST.end(), MST);

    //for(int i = 0; i < N-1; i++) {
    //    std::cout << "C++ (" << _MST[i].X << "," << _MST[i].Y << ") [" << _MST[i].W << "]" << std::endl;
    //}

    return MST;
}

long int father(long int i, std::vector<long int>& id) {
	//std::cout << "ID: [";
	//for(auto elem : id) {
	//	std::cout << elem << " ";
	//}
	//std::cout << "]" << std::endl;
	//std::cout << "i = " << i << std::endl;
    while(i != id[i]) {
        id[i] = id[id[i]];
        i = id[i];
		
		//std::cout << "i = " << i << std::endl;
    }
    return i;
}

void unite(long int p, long int q, std::vector<long int>& id) {
    long int i = father(p, id);
    long int j = father(q, id);
    id[i] = j;
}

void ParallelFK(
    std::vector<Edge> G,
    long int N,
    std::vector<Edge>& MST,
    std::vector<long int>& P, 
    std::unordered_map<long int, long int>& M,
    long int& CC
)
{
    if(G.size() < N) {
        Kruskal(G, N, MST, P, M, CC);
        return;
    }

    static std::default_random_engine rng;
    std::uniform_int_distribution<long int> uid(0, G.size()-1);

    long int p = uid(rng);

    std::vector<Edge> E_le, E_gt;
    // Do partition
    #pragma omp parallel
    {
        std::vector<Edge> local_E_le, local_E_gt;
        #pragma omp for nowait
        for(long int i = 0; i < G.size(); i++)
        {
            if(G[i].W <= G[p].W)
                local_E_le.push_back(G[i]);
            else
                local_E_gt.push_back(G[i]);
        }

        #pragma omp critical
        E_le.insert(E_le.end(), local_E_le.begin(), local_E_le.end());
        #pragma omp critical
        E_gt.insert(E_gt.end(), local_E_gt.begin(), local_E_gt.end());
    }
    
    ParallelFK(E_le, N, MST, P, M, CC);


    // Filter E_gt
    std::vector<Edge> filtered_E_gt;
    #pragma omp parallel
    {
        std::vector<Edge> local_E_gt;
        #pragma omp for nowait
        for(long int i = 0; i < E_gt.size(); i++)
        {
            auto e = E_gt[i];
            auto X_set = M.find(e.X);
            auto Y_set = M.find(e.Y);
            if(X_set == M.end() || Y_set == M.end() || father(M[e.X], P) != father(M[e.Y], P))
                local_E_gt.push_back(e);
        }

        #pragma omp critical
        filtered_E_gt.insert(filtered_E_gt.end(), local_E_gt.begin(), local_E_gt.end());
    }
    

    ParallelFK(filtered_E_gt, N, MST, P, M, CC);
}

void Kruskal(
    std::vector<Edge> G,
    long int N,
    std::vector<Edge>& MST,
    std::vector<long int>& P, 
    std::unordered_map<long int, long int>& M,
    long int& CC
)
{
    std::sort(std::execution::par_unseq, G.begin(), G.end(), [](Edge lhs, Edge rhs){return lhs.W < rhs.W;});
    
    for(auto e : G) {
        auto X_set = M.find(e.X);
        auto Y_set = M.find(e.Y);

        if(X_set == M.end() || Y_set == M.end()) {
            MST.push_back(e);

            if(X_set == M.end()) {
                M[e.X] = CC;
                P.push_back(CC);
                CC++;
        		Y_set = M.find(e.Y);
            }
            if(Y_set == M.end()) {
                M[e.Y] = CC;
                P.push_back(CC);
                CC++;
            }
            unite(M[e.X], M[e.Y], P);
			M[e.X] = P[M[e.X]];
			M[e.Y] = P[M[e.Y]];
        } else if(father(M[e.X], P) != father(M[e.Y], P) ) {
            MST.push_back(e);
            unite(M[e.X], M[e.Y], P);
			M[e.X] = P[M[e.X]];
			M[e.Y] = P[M[e.Y]];
        }
    }
}

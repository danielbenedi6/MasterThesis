#include <iostream>
#include <vector>
#include <random>
#include <cmath>
#include <fstream>
#include <sstream>
#include <chrono>
#include "parallel.hpp"

std::vector<Edge> generate(int N, double p, long int seed) {
	std::vector<Edge> G;
    std::mt19937 gen(seed); // Standard mersenne_twister_engine seeded
    std::uniform_real_distribution<> dis(0.0, 1.0);

	int v = 1;
	int w = -1;
	while(v < N) {
		double r = dis(gen);
		w = w + 1 + int(std::log(1-r)/std::log(1-p));
		while(w >= v && v < N) {
			w = w - v;
			v = v + 1;
		}

		if(v < N) {
			G.push_back(Edge{v,w, dis(gen)});
		}
	}

	return G;
}

/*
int main() {
    int N = 5000;
	double p = 0.75;
	long int seed = 1321412333121;

	auto graph = generate(N,p,seed);

    Edge *MST = ParallelFilterKruskal(graph.data(), graph.size(), N);

    for(int i = 0; i < N-1; i++)
        std::cout << MST[i].X << " " << MST[i].Y << " " << MST[i].W << std::endl;

    delete[] MST;
}*/

int main() {
	int repetitions = 10;
	std::ifstream file("./create_stats.csv");

	/*
	 1000,0.10,10,0,1.803130214s,1706466635555871667
	 1000,0.10,10,1,2.25787595s,1706466637359914999
	 1000,0.10,10,2,1.624787162s,1706466639617868493
	 1000,0.10,10,3,1.758576218s,1706466641242707038
	 1000,0.10,10,4,1.390425215s,1706466643001363946
	 1000,0.10,10,5,1.193007215s,1706466644391850580
	 1000,0.10,10,6,974.752936ms,1706466645584920276
	 1000,0.10,10,7,1.323148122s,1706466646559734098
	 1000,0.10,10,8,1.548497081s,1706466647882965467
	 1000,0.10,10,9,1.24271797s,1706466649431538670
	 */
	std::string line;
	while(std::getline(file, line) ){
			std::vector<std::string> fields;
			std::istringstream ss(line);
			while(std::getline(ss, line, ',')) {
				fields.push_back(line);
			}

			int N = std::stoi(fields[0]);
			double p = std::stod(fields[1]);
			long int seed = std::stol(fields[5]);
			
			auto graph = generate(N,p,seed);
			std::cout << N << "," << p;
			for(int exp = 0; exp < repetitions; exp++) {
				std::chrono::steady_clock::time_point begin = std::chrono::steady_clock::now();
				Edge *MST = ParallelFilterKruskal(graph.data(), graph.size(), N);
				std::chrono::steady_clock::time_point end = std::chrono::steady_clock::now();

				std::cout << "," << std::chrono::duration_cast<std::chrono::microseconds>(end - begin).count();
	
				delete[] MST;
			}
			std::cout << "\n";
	}
}

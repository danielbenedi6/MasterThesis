# Dynamic Pipeline Applied to Fully Dynamic Minimum Spanning Trees
---

This repository contains the code and materials for my master's thesis, which is an experimental analysis of the Dynamic Pipeline approach applied to the computation and maintenance of Minimum Spanning Trees (MSTs) in dynamic graphs.

A **Minimum Spanning Tree** is a subset of the edges of a connected, edge-weighted, undirected graph that connects all the vertices together without any cycles and with the minimum possible total edge weight.

The **Dynamic Pipeline** is a novel concurrent paradigm based on the Divide and Conquer approach. The key idea is to dynamically establish a set of filters that process input data and, when requested, return desired properties from the input data, such as the MST in this case.

## Repository Structure

The repository is organized as follows:

### `dp_mst/`
This is the main folder containing the codebase. Key components include:

- **Makefile**: Use this to build the main program.

- **go.mod**: Go Module defining file

- **internal/**: Contains various internal modules necessary for the project:
  - `mst/`: Code for computing an MST.
  - `dp/`: Code for building and communicating the pipeline.
  - `common/`: Utility functions and common tools needed across the project.

- **cmd/**: Different command-line interfaces to test various aspects of the project:
  - `cmdmst/`: Tests the behavior of `dp_kruskal`.
  - `dynfilterkruskal/`: Runs Kruskal's algorithm and Filter_Kruskal on dynamic graphs.
  - `filterkruskal/`: Generates static random graphs and tests Kruskal and Filter_Kruskal.
  - `inputtest/`: Tests different packages for input handling.
  - `instructions/`: Tests if using strings or numbers is more efficient for instructions.
  - `randomnumbers/`: Tests the quality of the random number generator.
  - `savegenerator/`: Generates random graphs, their distribution in the pipeline, and the corresponding savestate.

  Each of these binaries provides help information when run without arguments.

- **samples/**: Contains two examples of input data and how to load it into the program.

- **saves/**: Placeholder folder in which the savestate and loadstate operation is done.

### `DynGraphRepo/`
This is a dump of dynamic graphs obtained from [DynGraphRepo](https://dyngraphlab.github.io/#). It also contains a `stats.csv` with some properties of those graphs.

### `python_scripts/`
Here, there aare some utilities scripts such as:

- `filter_simulator.py`: Reads an input file and simulates how many and how dense the filters will be.

- `test_randoms.py`: Reads the random numbers generated and performs chi-square on buckets and analyses how many are suspicious or rejected.

- `generator.py` and `graph_generator.py`: Generators of inputs completly random input and input of random graph.

- `plot_*.py`: Different types of scripts to plot results.


## Usage

To build the main program, navigate to the `dp_mst` directory and run:

```bash
make
```

Ensure that you have Go installed and properly set up on your machine, as the project relies on Go's concurrency features.

To run the main program or any of the test interfaces, use the corresponding binary in the `cmd/` directory. Each binary can be run without arguments to display usage instructions and options.

For example, to test the behavior of `dp_kruskal`, you can navigate to the `cmdmst/` directory and run:

```bash
./bin/cmdmst
```

The program will display the available options and how to use them.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For any questions or issues, please feel free to reach out via email.

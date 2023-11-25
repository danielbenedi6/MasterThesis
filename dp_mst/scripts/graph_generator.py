import random
import networkx
from enum import IntEnum

class Operation(IntEnum):
    Insert  = 1
    Update  = 2
    Delete  = 3
    KMST    = 4
    EOF     = 5
    GraphOp = 6

def write_graph(G, name):
    with open(name, "w") as file:
        edges = list(G.edges(data=True))
        random.shuffle(edges)
        nodes = list(G.nodes())
        edge = 0

        while edge < len(edges):
            u,v,w = edges[edge]
            probability = random.random()
            if probability >= 0.5:
                u,v = v,u
            file.write(f"{Operation.Insert} {nodes.index(u)} {nodes.index(v)} {w['weight']}\n")
            edge += 1
        
        file.write(f"{Operation.GraphOp}\n")
        file.write(f"{Operation.KMST}\n")
        file.write(f"{Operation.EOF}\n")
            
        

def explode_nodes(G):
    nodes_to_explode = [node for node in G.nodes() if G.degree[node] > 3]

    for node in nodes_to_explode:
        neighbors = list(G.neighbors(node))
        num_neighbors = len(neighbors)
        new_nodes = []
        # Create d new nodes and connect them with zero-weight edges
        for i in range(num_neighbors):
            new_nodes += [f"{node}_{i}"]
            if len(new_nodes) > 1:
                G.add_edge(new_nodes[-2], new_nodes[-1], weight=0)  # Connect exploded nodes with weight 0
            G.add_edge(new_nodes[-1], neighbors[i], weight=G[node][neighbors[i]]["weight"])  # Connect to original neighbors 
        G.add_edge(new_nodes[0], new_nodes[-1], weight=0)
        G.remove_node(node)  # Remove the original node after exploding

def main():
    try:
        N = int(input("Enter the number of vertices: "))
        p = float(input("Enter the expected density: "))
        s = int(input("Enter the seed: "))

        random.seed(s)
        G = networkx.binomial_graph(N,p,seed=s)

        for u,v in G.edges():
            w = random.random() + 1
            G[u][v]["weight"] = w

        write_graph(G, f"input_test/{N}.requests")
        explode_nodes(G)
        write_graph(G, f"input_test/{N}_max3.requests")
    except ValueError:
        print("Invalid input.")
    except IOError:
        print("Error writing to the file.")

if __name__ == "__main__":
    main()

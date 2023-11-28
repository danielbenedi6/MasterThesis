import random
import networkx
import math
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
        random.shuffle(G)
        edge = 0

        while edge < len(G):
            u,v,w = G[edge]
            probability = random.random()
            if probability >= 0.5:
                u,v = v,u
            file.write(f"{Operation.Insert} {u} {v} {w}\n")
            edge += 1
        
        file.write(f"{Operation.GraphOp}\n")
        file.write(f"{Operation.KMST}\n")
        file.write(f"{Operation.EOF}\n")
            
def binomial_graph(n,p):
    G = []

    lp = math.log(1.0 - p)

    # Nodes in graph are from 0,n-1 (start with v as the second node index).
    v = 1
    u = -1
    while v < n:
        lr = math.log(1.0 - random.random())
        u = u + 1 + int(lr / lp)
        while u >= v and v < n:
            u = u - v
            v = v + 1
        if v < n:
            w = random.random() + 1 # generate random weight
            G += [(v, u, w)]
    return G        

def explode_nodes(G):
    N = max(G)[0] + 1
    deg_seq = [0 for _ in range(N)]
    for edge in G:
        u,v,_ = edge
        deg_seq[u] += 1
        deg_seq[v] += 1

    nodes_to_explode = [deg_seq[node] > 3 for node in range(N)]
    transform = dict()
    
    M = len(G)
    for edge in range(M):
        u,v,w = G[edge]

        if nodes_to_explode[u]:
            transform.setdefault(u, []).extend([N])
            u = N
            N += 1

        if nodes_to_explode[v]:
            transform.setdefault(v, []).extend([N])
            v = N
            N += 1
        
        G[edge] = (u,v,w)

    for original in transform:
        new_nodes = transform[original]
        for i in range(len(new_nodes)):
            G += [(new_nodes[i-1],new_nodes[i],0)]
    

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

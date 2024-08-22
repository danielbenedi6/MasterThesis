from math import exp, log
import random

def simulate_random_graph(N,p, normalize):
  filters = [0]
  roots = [0]
  mapToFilter = dict()
  v = 1
  u = -1
  while v < N:
    u += 1 + int( log(1 - random.random()) / log(1 - p))

    while u >= v and v < N:
      u = u - v
      v = v + 1
    if v < N: # Edge is in the graph and should be added to DP
      e = (u,v)
      # Normalization operation
      if normalize:
        if e[0] > e[1]:
            e[0], e[1] = e[1], e[0]
      else:
        # Random order of vertices
        if random.choice([True, False]):
            e[0], e[1] = e[1], e[0]
      
        if e[0] not in mapToFilter: # If vertex has not been assigned to a filter yet
            if roots[-1] == Fsize: # If last filter has reached its maximum, we should make a new one
                filters.append(0)
                roots.append(0)
            roots[-1] += 1
            mapToFilter[e[0]] = len(filters) - 1
        filters[mapToFilter[e[0]]] += 1

  return filters
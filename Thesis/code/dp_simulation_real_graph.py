from math import exp, log
import random

def simulate(path, normalize):
    filters = [dict()]
    mapToFilter = dict()

    file = open(path, "r")
    line = file.readline()

    V = set()
    E = set()

    while line:
        line = line.split()
        for i in range(3):
            line[i] = int(line[i])
        V.add(line[1])
        V.add(line[2])
        # Normalization operation
        if normalize:
          if line[1] > line[2]:
              line[1], line[2] = line[2], line[1]
        else:
          # Random order of vertices
          if random.choice([True, False]):
              line[1], line[2] = line[2], line[1]

        if line[0] == 1: # Insert Operation
            E.add((line[1],line[2]))
            if line[1] not in mapToFilter: # If vertex has not been assigned to a filter yet
                if len(filters[-1]) == Fsize: # If last filter has reached its maximum, we should make a new one
                    filters.append(dict())
                filters[-1][line[1]] = set()
                mapToFilter[line[1]] = len(filters) - 1
            filters[mapToFilter[line[1]]][line[1]].add(line[2])
        elif line[1] == 2:
            # Delete Operation
            if line[1] in mapToFilter and line[2] in filters[mapToFilter[line[1]]][line[1]]:
                filters[mapToFilter[line[1]]][line[1]].remove(line[2])
        line = file.readline()

    density = [ sum([ len(filter[root]) for root in filter ])  for filter in filters ]
    return density
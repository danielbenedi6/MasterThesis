import matplotlib.pyplot as plt
from math import exp

path = "movielens10m.requests"
Fsize = 10
t = 1

filters = [dict()]
mapToFilter = dict()

file = open(path, "r")
line = file.readline()

V = set()

while line:
    line = line.split()
    for i in range(3):
        line[i] = int(line[i])
    V.add(line[1])
    V.add(line[2])
    # Normalization operation
    if line[1] > line[2]:
        line[1], line[2] = line[2], line[1]
    if line[0] == 1: # Insert Operation
        if line[1] not in mapToFilter: # If vertex has not been assigned to a filter yet
            if len(filters[-1]) == Fsize: # If last filter has reached its maximum, we should make a new one
                filters.append(dict())
                Fsize = max(int(Fsize*1.01) , int(Fsize)+1)
                t += 1
            filters[-1][line[1]] = set()
            mapToFilter[line[1]] = len(filters) - 1
        filters[mapToFilter[line[1]]][line[1]].add(line[2])
    elif line[1] == 2:
        # Delete Operation
        if line[1] in mapToFilter and line[2] in filters[mapToFilter[line[1]]][line[1]]:
            filters[mapToFilter[line[1]]][line[1]].remove(line[2])
    line = file.readline()

density = [ sum([ len(filter[root]) for root in filter ])  for filter in filters ]
print(len(V))
plt.bar(list(range(len(density))), density)
plt.show()
import matplotlib.pyplot as plt
import numpy
import scipy.stats as stats

data = numpy.array(list(map(float, open("random_numbers.txt").readlines())))
data = data.reshape((int(len(data)/1000), 1000))

print("χ2 test")

k = int(len(data)/5) - 1
delta_b = 1.0 / k

bins = numpy.hstack(
                (
                    numpy.floor(data/delta_b), 
                    numpy.tile(numpy.arange(0,k),len(data)).reshape(len(data),k)
                )
            )
freq = numpy.array(list(map(lambda row: (numpy.unique(row, return_counts=True)[1] - 1), bins)))

test = list(map(lambda row : stats.chisquare(row), freq))

reject = 0
suspect = 0
almost = 0
alright = 0

for i in range(len(test)):
    if test[i].pvalue < 0.01 or test[i].pvalue > 0.99:
        reject += 1
    elif test[i].pvalue < 0.05 or test[i].pvalue > 0.95:
        suspect += 1
    elif test[i].pvalue < 0.1 or test[i].pvalue > 0.9:
        almost += 1
    else:
        alright += 1

print(f"Reject: {reject*100/len(test)}%. Suspect: {suspect*100/len(test)}%. Almost suspect: {almost*100/len(test)}%. Alright: {alright*100/len(test)}% ")
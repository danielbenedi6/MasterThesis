import seaborn as sns
import matplotlib.pyplot as plt
import pandas

filename = input("Enter the file path csv to analyze: ")

data = pandas.read_csv(filename, sep=";")
data["microseconds"] /= data["N"]

plt.figure(figsize=(15,10))
sns.lineplot(data, x="N", y="microseconds", hue="reader")
plt.xlabel("Number of lines")
plt.ylabel("Time per line [μs]")
#plt.yscale("log")
plt.title("Time comparison of reading file (N = 100)")
plt.savefig(filename+".png", dpi=900)

print(data[data["N"] >= 100000].groupby(["reader"]).mean())

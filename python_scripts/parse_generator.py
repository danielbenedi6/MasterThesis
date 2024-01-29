import seaborn as sns
import matplotlib.pyplot as plt
import pandas
import matplotlib
matplotlib.use("pgf")
matplotlib.rcParams.update({
    "pgf.texsystem": "pdflatex",
    'font.family': 'serif',
    'text.usetex': True,
    'pgf.rcfonts': False,
})

filename = input("Enter the file path csv to analyze: ")

data = pandas.read_csv(filename, sep=";")
data["microseconds"] /= data["N"]

plt.figure(figsize=(15,10))
sns.lineplot(data, x="N", y="microseconds", hue="reader")
plt.xlabel("Number of lines")
plt.ylabel("Time per line [$\mu$s]")
#plt.yscale("log")
plt.title("Time comparison of reading file")
plt.savefig("figures/" + filename.split(".")[0].split("/")[-1]+".pgf")


print(data[data["N"] >= 100000].groupby(["reader"]).mean())

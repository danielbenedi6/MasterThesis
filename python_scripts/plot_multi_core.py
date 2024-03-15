import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd
from math import log, sqrt

def plot(data):
    custom_palette = sns.color_palette("husl", n_colors=len(data['p'].unique()))
    sns.lineplot(x="N", y="Mean", hue="p", style="Type", data=data, markers=True, palette=custom_palette)

    plt.xlabel("Number of nodes")
    plt.ylabel("Execution Time [µs]")

    plt.yscale("log")
    plt.xscale("log")

def plot_no_p(data):
    custom_palette = sns.color_palette("husl", n_colors=len(data['Type'].unique()))
    sns.lineplot(x="N", y="Mean", hue="Type", data=data, markers=True, palette=custom_palette)

    plt.xlabel("Number of nodes")
    plt.ylabel("Execution Time [µs]")

    plt.yscale("log")
    plt.xscale("log")

def plot_cores(data):
    custom_palette = sns.color_palette("husl", n_colors=len(data['Type'].unique()))
    sns.lineplot(x="Cores", y="Mean", hue="Type", data=data, markers=True, palette=custom_palette)

    plt.xlabel("Number of cores")
    plt.ylabel("Execution Time [µs]")

    plt.yscale("log")
    plt.xscale("log", base=2)

data = pd.DataFrame(columns=["Type", "N", "p", "Mean","Cores"])
for core in [1,2,4,8,16]:
    df = pd.read_csv(f"core{core}.csv", names=["N", "p", "Fsize", "Try"] + [f"T{i}" for i in range(11)])
    for i, row in df.iterrows():
        if row["Fsize"] == int(sqrt(row["N"])):
            df.loc[i, "Type"] = "DP-Kruskal-sqrt"
        elif row["Fsize"] == int(log(row["N"])):
            df.loc[i, "Type"] = "DP-Kruskal-log"
        else:
            df.loc[i, "Type"] = "DP-Kruskal-const"

        for t in [f"T{i}" for i in range(11)]:
            if "µs" in row[t]:
                df.loc[i,t] = float(row[t][:-2])
            elif "ms" in row[t]:
                df.loc[i,t] = float(row[t][:-2])*1e3
            elif "m" in row[t]:
                df.loc[i,t] = (float(row[t].split("m")[0])*60+float(row[t].split("m")[1][:-1]))*1e6
            else:
                df.loc[i,t] = float(row[t][:-1])*1e6
    df["Mean"] = df[["T"+str(i+1) for i in range(10)]].sum(axis=1) / 10.0
    df["Cores"] = core
    data = pd.concat([data, df[["Type", "N", "p", "Mean","Cores"]]])
    
    # Plot all densities
    plt.figure(figsize=(16,16))
    plot(df)
    plt.title(f"Evolution with random graph in {core}-core")
    plt.savefig(f"AllDensities-Core{core}.pdf", format="pdf")
    plt.show()
    
    # Plot some densities individually
    plt.figure(figsize=(16,6))
    ax1 = plt.subplot(1,3,1)
    plot_no_p(df[df["p"] == 0.1])
    plt.title(f"Evolution with random graph of expected density of 0.1 in {core}-core", wrap=True)
    plt.subplot(1,3,2, sharey=ax1)
    plot_no_p(df[df["p"] == 0.5])
    plt.title(f"Evolution with random graph of expected density of 0.5 in {core}-core", wrap=True)
    plt.subplot(1,3,3, sharey=ax1)
    plot_no_p(df[df["p"] == 0.9])
    plt.title(f"Evolution with random graph of expected density of 0.9 in {core}-core", wrap=True)
    plt.savefig(f"SomeDensities-Core{core}.pdf", format="pdf")
    plt.show()

plt.figure(figsize=(32,32))
i = 1
for N in [1000, 5000, 10000]:
    for p in [0.1,0.5,0.9]:
        plt.subplot(3,3,i)
        plot_cores(data[ (data["p"] == p) & (data["N"] == N) ] )
        plt.title(f"Evolution with G(n={N}, p={p})")
        i += 1

plt.savefig("CoresComparison.pdf", format="pdf")
plt.show()


global_mean = data.groupby(["Type", "Cores", "N", "p"])["Mean"].mean().reset_index()
global_mean["Speed-Up"] = 1
global_mean["Utilization"] = 1
global_mean["Expected edges"] = 1
for i, row in global_mean.iterrows():
    global_mean.loc[i, "Speed-Up"] = global_mean[ (global_mean["Type"] == row["Type"]) & (global_mean["N"] == row["N"]) & (global_mean["p"] == row["p"]) & (global_mean["Cores"] == 1) ]["Mean"].to_numpy()[0] / row["Mean"]
    global_mean.loc[i, "Utilization"] = global_mean.loc[i, "Speed-Up"] / row["Cores"]
    global_mean.loc[i, "Expected edges"] = row["N"] * (row["N"] - 1) / 2 * row["p"]

sns.lineplot(global_mean, x="Expected edges", y="Speed-Up", hue="Type", style="Cores")
plt.show()

ax1 = plt.subplot(1,3,1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-const"], x="Expected edges", y="Speed-Up", hue="Cores")
plt.title("Speed-up of DP-Kruskal-const")
plt.subplot(1,3,2, sharey=ax1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-log"], x="Expected edges", y="Speed-Up", hue="Cores")
plt.title("Speed-up of DP-Kruskal-log")
plt.subplot(1,3,3, sharey=ax1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-sqrt"], x="Expected edges", y="Speed-Up", hue="Cores")
plt.title("Speed-up of DP-Kruskal-sqrt")

plt.show()

ax1 = plt.subplot(1,3,1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-const"], x="Expected edges", y="Utilization", hue="Cores")
plt.title("Efficiency of DP-Kruskal-const")
plt.subplot(1,3,2, sharey=ax1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-log"], x="Expected edges", y="Utilization", hue="Cores")
plt.title("Efficiency of DP-Kruskal-log")
plt.subplot(1,3,3, sharey=ax1)
sns.lineplot(global_mean[global_mean["Type"] == "DP-Kruskal-sqrt"], x="Expected edges", y="Utilization", hue="Cores")
plt.title("Efficiency of DP-Kruskal-sqrt")

plt.show()


plt.figure(figsize=(32,32))
i = 1
for N in [1000, 5000, 10000]:
    for p in [0.1,0.5,0.9]:
        plt.subplot(3,3,i)
        sns.lineplot(global_mean[ (global_mean["p"] == p) & (global_mean["N"] == N)], x="Cores", y="Utilization", hue="Type")
        plt.xscale("log", base=2)
        i += 1
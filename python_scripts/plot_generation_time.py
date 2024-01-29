import numpy as np
import pandas as pd
import seaborn as sns
import matplotlib.pyplot as plt
import matplotlib
matplotlib.use("pgf")
matplotlib.rcParams.update({
    "pgf.texsystem": "pdflatex",
    'font.family': 'serif',
    'text.usetex': True,
    'pgf.rcfonts': False,
})

df = pd.read_csv(
        "create_stats.csv", 
        names=["N", "edge probability", "Fsize", "id", "time", "seed"], 
        converters={"time": pd.Timedelta}
    )

df['time [s]'] = df['time'].astype('timedelta64[us]')/1000000
df["type"] = "fixed"
df.loc[df["Fsize"] == np.log(df["N"]).astype('int'), "type"] = "log"
df.loc[df["Fsize"] == np.sqrt(df["N"]).astype('int'), "type"] = "sqrt"

sns.lineplot(x="N", 
             y="time [s]", 
             hue="edge probability", 
             data=df[df["type"] == "fixed"]
            ).set_title("Construction time for fixed size of filters")
plt.xscale("log")
plt.yscale("log")
plt.savefig("figures/gen_time_fixed_size.pgf")

sns.lineplot(x="N", 
             y="time [s]", 
             hue="edge probability", 
             data=df[df["type"] == "sqrt"]
            ).set_title("Construction time for size of filters of $\sqrt{N}$")
plt.xscale("log")
plt.yscale("log")
plt.savefig("figures/gen_time_sqrt_size.pgf")

sns.lineplot(x="N", 
             y="time [s]", 
             hue="edge probability", 
             data=df[df["type"] == "log"]
            ).set_title("Construction time for size of filters of $\log{N}$")
plt.xscale("log")
plt.yscale("log")
plt.savefig("figures/gen_time_log_size.pgf")

import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D

df = pd.read_csv("benchmark_results.txt")

grouped = df.groupby(["degree", "rows"])["search_ns"].mean().reset_index()

degrees = sorted(df["degree"].unique())
rows = sorted(df["rows"].unique())

X, Z = np.meshgrid(degrees, rows)
Y = np.zeros_like(X, dtype=float)

for i, r in enumerate(rows):
    for j, d in enumerate(degrees):
        val = grouped[(grouped["degree"] == d) & (grouped["rows"] == r)]["search_ns"].values
        if len(val) > 0:
            Y[i, j] = val[0]

fig = plt.figure()
ax = fig.add_subplot(111, projection="3d")

# log scale to spread out differences
Y_log = np.log1p(Y)
norm = plt.Normalize(Y_log.min(), Y_log.max())
colors = plt.cm.plasma(norm(Y_log))

surf = ax.plot_surface(X, Z, Y, facecolors=colors)

mappable = plt.cm.ScalarMappable(cmap="plasma", norm=norm)
mappable.set_array(Y_log)
fig.colorbar(mappable, ax=ax, shrink=0.5, label="search time (ns, log scale)")

ax.set_xlabel("degree")
ax.set_zlabel("search time (ns)")
ax.set_ylabel("rows")
ax.set_title("B-tree search performance")

print(df)
plt.show()

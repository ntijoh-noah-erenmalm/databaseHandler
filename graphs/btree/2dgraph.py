import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

df = pd.read_csv("benchmark_results.txt")
df["search_us"] = df["search_ns"] / 1000
mean_df = df.groupby(["degree", "rows"])["search_us"].mean().reset_index()
degrees = sorted(df["degree"].unique())
rows = sorted(df["rows"].unique())

fig, ax = plt.subplots(figsize=(12, 6))

# generate a unique color for each degree
colors = plt.cm.tab20(np.linspace(0, 1, len(degrees)))

for i, degree in enumerate(degrees):
    subset = mean_df[mean_df["degree"] == degree]
    ax.plot(subset["rows"], subset["search_us"], marker="o", label=f"degree {degree}", color=colors[i])

ax.set_xlabel("rows")
ax.set_ylabel("search time (µs)")
ax.set_title("B-tree search time by row count and degree")
ax.legend()
ax.set_xticks(rows)
ax.set_xticklabels([f"{r:,}" for r in rows], rotation=45)

plt.tight_layout()
plt.show()

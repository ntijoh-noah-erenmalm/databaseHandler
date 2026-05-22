import pandas as pd
import matplotlib.pyplot as plt

df = pd.read_csv("linear_benchmark_results.txt")

df["search_us"] = df["search_ns"] / 1000

def get_target_type(row):
    if row["search_target"] == 0:
        return "first"
    elif row["search_target"] == row["rows"] - 1:
        return "last"
    else:
        return "random"

df["target_type"] = df.apply(get_target_type, axis=1)

fig, ax = plt.subplots(figsize=(10, 6))

for target in ["first", "last", "random"]:
    subset = df[df["target_type"] == target]
    ax.plot(subset["rows"], subset["search_us"], marker="o", label=target)

mean_df = df.groupby("rows")["search_us"].mean().reset_index()
ax.plot(mean_df["rows"], mean_df["search_us"], marker="o", linestyle="--", color="black", label="mean")

ax.set_xlabel("rows")
ax.set_ylabel("search time (µs)")
ax.set_title("linear search time by row count")
ax.legend()

# use actual row values as ticks instead of log scale
rows = sorted(df["rows"].unique())
ax.set_xticks(rows)
ax.set_xticklabels([f"{r:,}" for r in rows], rotation=45)

plt.tight_layout()
plt.show()

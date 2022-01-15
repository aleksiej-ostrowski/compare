# -------------------------------- #
#                                  #
#  version 0.0.1                   #
#                                  #
#  Aleksiej Ostrowski, 2022        #
#                                  #
#  https://aleksiej.com            #
#                                  #
# -------------------------------- #

import matplotlib.pyplot as plt
import random
import time
import json
from types import SimpleNamespace
import math
import sys

# import humanize
# import datetime as dt


# https://stackoverflow.com/questions/37765197/darken-or-lighten-a-color-in-matplotlib
def adjust_lightness(color, amount=0.5):
    import matplotlib.colors as mc
    import colorsys

    try:
        c = mc.cnames[color]
    except Exception:
        c = color
    c = colorsys.rgb_to_hls(*mc.to_rgb(c))
    return colorsys.hls_to_rgb(c[0], max(0, min(1, amount * c[1])), c[2])


colors = [
    "gray",
    "rosybrown",
    "darksalmon",
    "sandybrown",
    "tan",
    "khaki",
    "darkkhaki",
    "lightseagreen",
    "darkcyan",
    "darkturquoise",
    "deepskyblue",
    "mediumpurple",
    "plum",
    "palevioletred",
    "darkseagreen",
    "lightskyblue",
    "lightsteelblue",
]


markers = ["o", "v", "^", "<", ">", "s", "p", "P", "*", "h", "X", "D"]
linestyles = ["solid", "dotted", "dashed", "dashdot"]

random.seed(time.time())

data = sys.stdin.read()

# print(data)
# exit()

p = json.loads(data, object_hook=lambda d: SimpleNamespace(**d))

"""
{
"data":  [
        [200, 1000, 10, 2100, 2501, 3501],
        [505, 5010, 50, 5200, 5502, 5502],
        [210, 1020, 30, 2300, 2503, 3503],
        [220, 1030, 40, 2400, 2504, 3504],
        [225, 1040, 50, 2500, 2505, 3505]
        ],
"labels": ["1", "2", "3", "4", "5"],
"x": [10000, 100000, 500000, 1000000, 5000000, 10000000],
"xlabel": "N",
"ylabel": "Time, sec.",
"title": "Algorithm's compare"
}

"""

len_data = len(p.Data)

n_len = float(len_data) / len(colors)

if n_len > 1.0:
    colors *= math.ceil(n_len)

n_len = float(len_data) / len(markers)

if n_len > 1.0:
    markers *= math.ceil(n_len)

n_len = float(len_data) / len(linestyles)

if n_len > 1.0:
    linestyles *= math.ceil(n_len)

random.shuffle(colors)
random.shuffle(markers)
random.shuffle(linestyles)

# ns -> ms
for i, r in enumerate(p.Data):
    for j, x in enumerate(r):
        p.Data[i][j] *= 1e-6


for i, e in enumerate(p.Data):
    plt.plot(
        p.X,
        e,
        linestyle=linestyles[i],
        marker=markers[i],
        label=p.Labels[i],
        color=colors[i],
    )

    for_filter = zip(p.X, e)

    only_this = filter(
        eval(
            "lambda x: {r}".format(
                r=p.Xfilter.replace("x", "x[0]").replace("y", "x[1]")
            )
        ),
        for_filter,
    )

    for a, b in only_this:
        """
        b_fmt = humanize.precisedelta(
            dt.timedelta(milliseconds=b), minimum_unit="milliseconds"
        )
        b_fmt = humanize.scientific(b, precision=1)
        """
        b_fmt = str(round(b, 1))
        plt.text(a, b + 2.0, b_fmt, color=adjust_lightness(colors[i]))

plt.grid(axis="x")  # , alpha = 0.15)
plt.grid(axis="y")  # , alpha = 0.15)

plt.xlabel(p.Xlabel)
plt.ylabel(p.Ylabel + ", ms")
plt.title(p.Title)

plt.legend()
plt.show()

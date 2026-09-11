"""Print the balance baseline from a simulate run.

The loop's brief says to parse the JSON rather than grep it, and every tick was
doing that by hand with a here-document. This is the same three lines, kept
where they can be run rather than retyped.
"""

import json
import sys

# The magpie is last because it is not a strategy anybody would follow. It has
# no plan: it works the docks until it can afford to be curious and then takes
# whatever it has taken least. It is here so that the parts of the game the
# eight policies with plans never touch are priced at all — they took 78 of the
# game's 116 kinds of action, and everything else went unmeasured.
ORDER = (
    "worker",
    "investor",
    "defiant",
    "reckless",
    "thief",
    "smuggler",
    "racketeer",
    "publican",
    "magpie",
    "distiller",
)

# Days as well as deaths and cash.
#
# The default campaign was two hundred commands, which is about twelve game
# days, and the simulator warns on every run that its own city measures need
# twenty. Everything this project believed about its economy came out of
# campaigns too short to see it: at two hundred commands a publican ends nearly
# the poorest policy in the game and at four hundred it ends the richest, at
# $27,516 against the worker's $23,252. Running what you hold does pay. It pays
# later than anybody here had ever looked.
# And how much of the cash figure is the policy rather than the dice.
#
# This line was printed as a bare median for a long time and read as a ranking
# every time. It is not one. A campaign's final cash has a spread of about six
# thousand dollars around a median of twenty-six, so two policies a thousand
# apart are the same policy as far as a hundred runs can say — and the night the
# distiller was written this log called it "the richest policy in the harness"
# on a gap of 979 against a standard error of about 800.
#
# So the figure carries its own error bar, and any two policies whose bars
# overlap are named as what they are: not separable by this measure. The error
# on a median is about 1.25 times the standard deviation over the root of the
# count, which is close enough for deciding whether to believe a gap. Checked
# against the run that prompted it: a hundred distiller campaigns have a
# standard deviation of $6,371, which gives $796, and the printed bar says 796.
def spread(values):
    n = len(values)
    if n < 2:
        return 0
    mean = sum(values) / n
    sd = (sum((v - mean) ** 2 for v in values) / n) ** 0.5
    return int(1.25 * sd / (n ** 0.5))


run = json.load(open(sys.argv[1]))
summary = run["summary"]
cash = {}
for campaign in run.get("campaigns", []):
    cash.setdefault(campaign["strategy"], []).append(campaign["cash"])

played = [k for k in ORDER if k in summary]
error = {k: spread(cash.get(k, [])) for k in played}
print("deaths", "/".join(str(summary[k]["deaths"]) for k in played))
print("cash", "/".join(
    "%d±%d" % (summary[k]["median_final_cash"], error[k]) for k in played))
print("days", "/".join(str(summary[k]["median_game_days"]) for k in played))

# Which of them this run cannot tell apart. Read before believing any ordering
# in the line above.
same = []
for i, a in enumerate(played):
    for b in played[i + 1:]:
        gap = abs(summary[a]["median_final_cash"] - summary[b]["median_final_cash"])
        if gap <= error[a] + error[b]:
            same.append("%s~%s" % (a, b))
if same:
    print("not separable by cash:", " ".join(same))

city = run.get("city_alone", {}).get("totals")
if city:
    print("city_alone", city)

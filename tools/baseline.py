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
)

summary = json.load(open(sys.argv[1]))["summary"]
print("deaths", "/".join(str(summary[k]["deaths"]) for k in ORDER if k in summary))
print("cash", "/".join(str(summary[k]["median_final_cash"]) for k in ORDER if k in summary))
city = json.load(open(sys.argv[1])).get("city_alone", {}).get("totals")
if city:
    print("city_alone", city)

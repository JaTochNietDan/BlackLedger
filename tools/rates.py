"""What every action in this game pays for the time it takes.

The project could not answer that. The balance baseline reports where a policy
ends up, which is a fact about the policy; nothing measured what a single thing
is worth. Run a traced simulation and this reads the price of everything out of
it:

    go run ./cmd/simulate -strategies magpie -runs 120 -steps 2500 -trace \\
        > .runtime/trace.json
    python3 tools/rates.py .runtime/trace.json

Three things it has to do or it says nothing true.

**Per minute, not per command.** Time is the real currency here — the clock is
what a player spends. An action taking ninety minutes and paying $80 is worse
than one taking five and paying $20.

**Net of doing nothing.** The city pays the player while they stand there, so
every action looks profitable if you only take the cash difference. Waiting is
the baseline and everything is reported over it.

**With enough samples.** The exploring policy takes each kind of thing about
once a run, so twelve runs is twelve samples and the table is noise. Forty is
the floor here and a hundred runs is what makes it steady. The first version of
this measurement was read off twelve runs and its top row was a gambling
outcome.

**It cannot see deferred money, and this matters more than it sounds.** The
measure attributes a cash change to whatever command was running when it
landed, and several things in this game pay later than they are decided. Two
that will mislead anybody reading the table cold:

A bet and its settlement are separate commands, so `play` reads as the worst
rate in the game and `stand` as the best. They net out — across eight hundred
and twenty commands at the tables and machines the house keeps about sixty cents
a command net of ambient income, which is break-even within variance.

Sending somebody on collections pays $65 two hours after the decision, so
`delegate` reads as costing six cents a minute and doing nothing. A publican
sends 140 rounds a campaign; $109,265 of collections across twelve campaigns
landed while the player was asleep and read as up to 35% of what resting
appeared to pay. "Resting is the publican's economy" is the wrong sentence to
take away from this table, and it is the one the table says.

So: an action that decides something now and pays for it later reads as free.
Check what a command actually does before believing its row.
"""

import collections
import json
import statistics
import sys

FLOOR = 40


def main():
    runs = json.load(open(sys.argv[1]))["campaigns"]
    cash = collections.defaultdict(list)
    mins = collections.defaultdict(list)
    for r in runs:
        trace = r.get("trace") or []
        for a, b in zip(trace, trace[1:]):
            kind = (a["command"].get("kind") or "").split(":")[0]
            minutes = b["minute"] - a["minute"]
            if minutes <= 0:
                continue
            cash[kind].append(b["cash"] - a["cash"])
            mins[kind].append(minutes)
    if not mins.get("wait"):
        sys.exit("no waiting in this run, so there is no baseline to measure against")
    doing_nothing = sum(cash["wait"]) / sum(mins["wait"])

    rows = []
    for kind in cash:
        if len(cash[kind]) < FLOOR:
            continue
        rate = sum(cash[kind]) / sum(mins[kind]) - doing_nothing
        rows.append((rate, len(cash[kind]), statistics.median(mins[kind]), kind))
    rows.sort(reverse=True)

    print(f"{len(runs)} campaigns. Doing nothing pays {doing_nothing:.2f} a minute.")
    print(f"{len(rows)} actions taken {FLOOR} times or more.\n")
    print(f"{'$/min over':>11} {'times':>6} {'mins':>5}  action")
    for rate, times, minutes, kind in rows:
        print(f"{rate:11.2f} {times:6d} {minutes:5.0f}  {kind}")
    thin = sorted(k for k in cash if len(cash[k]) < FLOOR)
    if thin:
        print(f"\ntoo few samples to say anything about: {', '.join(thin)}")


if __name__ == "__main__":
    main()

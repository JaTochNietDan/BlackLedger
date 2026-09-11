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

**Money that arrives late.** A cash change lands on whatever command was running
when it turned up, which is not the command that sent for it. Sending somebody
on collections pays $65 two hours after the decision, so `delegate` read as
costing six cents a minute and doing nothing, and whichever command the crew
came back during read as generous — for a policy that rests four hours at a
time, that is resting. A third of what a publican appeared to earn by sitting
still was collections landing while it was asleep.

The trace records what came home and what it paid, so this puts it back against
the decision. The difference is not a correction, it is the opposite answer:
`delegate` is the best rate a publican has at $3.24 a minute, because the round
costs the crew two hours and the player fifteen minutes.

What is still unattributed is the felt. A bet and its settlement are separate
commands, so `play` reads as the worst rate in the game and `stand` as the best.
They net out — across eight hundred and twenty commands at the tables the house
keeps about sixty cents a command net of ambient income, break-even within
variance — but the two halves are not joined up here. Anything else that decides
now and pays later will read as free until the trace says otherwise.
"""

import collections
import json
import statistics
import sys

FLOOR = 40

# The action that sends for money which arrives later. The trace records when a
# collection came home and what it paid, so it can be taken off the command it
# landed on and given to the one that sent for it.
DEFERRED = "delegate"


def main():
    runs = json.load(open(sys.argv[1]))["campaigns"]
    cash = collections.defaultdict(list)
    mins = collections.defaultdict(list)
    late = collections.Counter()
    for r in runs:
        trace = r.get("trace") or []
        for a, b in zip(trace, trace[1:]):
            kind = (a["command"].get("kind") or "").split(":")[0]
            minutes = b["minute"] - a["minute"]
            if minutes <= 0:
                continue
            # Money that arrived late during this command belongs to whatever
            # sent for it, not to whatever was running when it turned up.
            arrived = a.get("settled", 0) * a.get("per_task", 0)
            cash[kind].append(b["cash"] - a["cash"] - arrived)
            mins[kind].append(minutes)
            late[DEFERRED] += arrived
    for kind, amount in late.items():
        if amount and mins.get(kind):
            cash[kind].append(amount)
    # Standing still, however this policy stands still. Waiting is the purest
    # form of it; a policy that never waits but rests is doing the same thing
    # as far as money is concerned, and a publican never waits at all.
    idle = "wait" if mins.get("wait") else "rest"
    if not mins.get(idle):
        sys.exit("this policy never stands still, so there is no baseline to measure against")
    doing_nothing = sum(cash[idle]) / sum(mins[idle])

    rows = []
    for kind in cash:
        if len(cash[kind]) < FLOOR:
            continue
        rate = sum(cash[kind]) / sum(mins[kind]) - doing_nothing
        rows.append((rate, len(cash[kind]), statistics.median(mins[kind]), kind))
    rows.sort(reverse=True)

    print(f"{len(runs)} campaigns. Standing still ({idle}) pays {doing_nothing:.2f} a minute.")
    print(f"{len(rows)} actions taken {FLOOR} times or more.\n")
    print(f"{'$/min over':>11} {'times':>6} {'mins':>5}  action")
    for rate, times, minutes, kind in rows:
        print(f"{rate:11.2f} {times:6d} {minutes:5.0f}  {kind}")
    thin = sorted(k for k in cash if len(cash[k]) < FLOOR)
    if thin:
        print(f"\ntoo few samples to say anything about: {', '.join(thin)}")


if __name__ == "__main__":
    main()

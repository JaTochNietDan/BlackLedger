"""Find, and prove the removal of, declarations the stylesheet already overrides.

`src/style.css` has grown by appending. A block gets rebuilt, the new version
goes on the end, and the old one stays where it was. Where two rules carry the
identical selector the later one simply wins, so the earlier declarations are
dead — harmless until somebody reads the earlier block, believes it, and edits
it. That is not hypothetical: the drums did not spin because an older `.drum`
rule set a centred grid and the newer rule never said otherwise.

Two commands, and the second is the one that matters.

    python3 tools/cssdead.py list
        Name every declaration a later rule with the same selector overrides.
        The test in cmd/blackledger/sources_test.go pins this count.

    python3 tools/cssdead.py same <before.css>
        Resolve every property of every selector to its last-wins value, in the
        old file and the new one, and print what differs. Unpicking a block is
        only safe if this prints nothing but the drops you meant.

Cleaning up is done a block at a time, by hand, reading both copies. It was
once attempted as a bulk regex pass over all fifty and reverted twice: comments
have to be blanked with spaces rather than cut or every byte offset past the
first one is wrong, and this file is not in the gate's Prettier set, so running
a formatter over it rewrites thousands of lines that have nothing to do with
the change.
"""

import collections
import re
import sys

STYLESHEET = "src/style.css"


def rules(path):
    """Every top-level rule as (selector, body, line). Comments are blanked
    with spaces rather than removed so offsets and line numbers still line up
    with the file as written."""
    raw = open(path).read()
    src = re.sub(
        r"/\*.*?\*/", lambda m: re.sub(r"[^\n]", " ", m.group(0)), raw, flags=re.S
    )
    out, depth, buf, line = [], 0, "", 1
    for i, c in enumerate(src):
        if c == "\n":
            line += 1
        if c == "{":
            if depth == 0:
                sel, body_at, at = buf.strip(), i + 1, line
            depth += 1
            buf = ""
        elif c == "}":
            depth -= 1
            if depth == 0:
                out.append((sel, src[body_at:i], at))
                buf = ""
        else:
            buf += c
    return out


def declared(body):
    return [
        (m.group(1), " ".join(m.group(2).split()))
        for m in re.finditer(r"([a-z-]+)\s*:\s*([^;]*)", body)
    ]


def dead(path):
    """Declarations a later rule with the same selector overrides. An earlier
    !important still wins, so it is not dead and neither is what it beats."""
    rs = rules(path)
    where = collections.defaultdict(list)
    for i, (sel, _, _) in enumerate(rs):
        if sel.startswith("@"):
            continue
        for one in sel.split(","):
            where[one.strip()].append(i)
    found = []
    for idxs in where.values():
        if len(idxs) < 2:
            continue
        for pos, idx in enumerate(idxs[:-1]):
            later = {
                prop
                for j in idxs[pos + 1 :]
                for prop, val in declared(rs[j][1])
                if "!important" not in val
            }
            for prop, val in declared(rs[idx][1]):
                if prop in later and "!important" not in val:
                    found.append((rs[idx][2], rs[idx][0], prop, val))
    return sorted(found)


def resolved(path):
    """What each selector's properties actually come out as: last one wins."""
    out = collections.defaultdict(dict)
    for sel, body, _ in rules(path):
        if sel.startswith("@"):
            continue
        for one in sel.split(","):
            for prop, val in declared(body):
                out[one.strip()][prop] = val
    return out


def main():
    what = sys.argv[1] if len(sys.argv) > 1 else "list"
    if what == "list":
        found = dead(STYLESHEET)
        for line, sel, prop, val in found:
            print(f"line {line:5d}  {sel.split(',')[0][:30]:32s} {prop}: {val[:40]}")
        print(len(found), "dead declarations")
        return
    if what == "same":
        before, after = resolved(sys.argv[2]), resolved(STYLESHEET)
        moved = 0
        for sel in sorted(set(before) | set(after)):
            a, b = before.get(sel, {}), after.get(sel, {})
            for prop in sorted(set(a) | set(b)):
                if a.get(prop) != b.get(prop):
                    moved += 1
                    print(f"  {sel} {prop}: {a.get(prop)!r} -> {b.get(prop)!r}")
        print(moved, "properties resolve differently")
        return
    sys.exit(__doc__)


if __name__ == "__main__":
    main()

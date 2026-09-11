"""Find, remove and prove the removal of declarations the stylesheet overrides.

`src/style.css` has grown by appending. A block gets rebuilt, the new version
goes on the end, and the old one stays where it was. Where two rules carry the
identical selector the later one simply wins, so the earlier declarations are
dead — harmless until somebody reads the earlier block, believes it, and edits
it. That is not hypothetical: the drums did not spin because an older `.drum`
rule set a centred grid and the newer rule never said otherwise.

    python3 tools/cssdead.py list
        Name every dead declaration. The test in
        cmd/blackledger/sources_test.go pins this count.

    python3 tools/cssdead.py strip
        Remove them, and any rule left with nothing in it. Nothing else in the
        file is touched: no reflowing, no tidying of semicolons elsewhere, no
        formatter — this stylesheet is not in the gate's Prettier set and a
        formatter run rewrites thousands of unrelated lines.

    python3 tools/cssdead.py same <before.css>
        Resolve every property of every selector to its last-wins value in both
        files and print what differs. A cleanup is only safe if this prints
        nothing but the drops that were meant. Run it every time.

Three things this had to get right, each of them wrong first and each caught by
`same` or by a shredded file:

Comments are blanked with spaces rather than cut, so every byte offset still
describes the file as written. Deleting them shifts everything past the first
comment and the splice lands mid-word.

A declaration is dead only when EVERY selector on its rule is later re-set for
that property. A rule can carry several — this file has
`.city-view-switch button,.street-inspect button,.street-journey button` in one
— and cutting because one of the three is shadowed takes the declaration from
the other two as well. That shipped once and `same` reported four properties
resolving differently on a change meant to change nothing.

One function finds them and both commands use it. The first version of `strip`
carried its own copy of the search, which is the oldest fault in this project:
two implementations of one rule that drift apart the first time either is
touched.
"""

import collections
import re
import sys

STYLESHEET = "src/style.css"


def blanked(raw):
    """The file with comments turned to spaces, so offsets still line up."""
    return re.sub(
        r"/\*.*?\*/", lambda m: re.sub(r"[^\n]", " ", m.group(0)), raw, flags=re.S
    )


def rules(raw):
    """Every top-level rule as (selector, body start, body end, line)."""
    src = blanked(raw)
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
                out.append((sel, body_at, i, at))
                buf = ""
        else:
            buf += c
    return out


DECLARATION = re.compile(r"([a-z-]+)\s*:\s*([^;]*)")


def declared(src, a, b):
    """Declarations in a body, with where each one sits in the file."""
    return [
        (m.group(1), " ".join(m.group(2).split()), a + m.start(), a + m.end())
        for m in DECLARATION.finditer(src[a:b])
    ]


def dead(raw):
    """Dead declarations as (line, selector, property, value, start, end).

    An earlier !important still wins, so it is not dead and neither is what it
    would have beaten.
    """
    src = blanked(raw)
    rs = rules(raw)
    where = collections.defaultdict(list)
    for i, (sel, _, _, _) in enumerate(rs):
        if not sel.startswith("@"):
            for one in sel.split(","):
                where[one.strip()].append(i)

    def overridden(idx, prop):
        for one in sel_of(rs, idx):
            if not any(
                later == prop and "!important" not in val
                for j in where[one]
                if j > idx
                for later, val, _, _ in declared(src, rs[j][1], rs[j][2])
            ):
                return False
        return True

    found = []
    for idx, (sel, a, b, line) in enumerate(rs):
        if sel.startswith("@"):
            continue
        for prop, val, at, end in declared(src, a, b):
            if "!important" not in val and overridden(idx, prop):
                found.append((line, sel, prop, val, at, end))
    return sorted(found)


def sel_of(rs, idx):
    return [one.strip() for one in rs[idx][0].split(",")]


def resolved(raw):
    """What each selector's properties actually come out as: last one wins."""
    src = blanked(raw)
    out = collections.defaultdict(dict)
    for sel, a, b, _ in rules(raw):
        if sel.startswith("@"):
            continue
        for one in sel.split(","):
            for prop, val, _, _ in declared(src, a, b):
                out[one.strip()][prop] = val
    return out


def strip(raw):
    """The file with every dead declaration gone, and every emptied rule."""
    out = raw
    for _, _, _, _, a, b in sorted(dead(raw), key=lambda d: d[4], reverse=True):
        j = b
        while j < len(out) and out[j] in " \t":
            j += 1
        if j < len(out) and out[j] == ";":
            j += 1
        out = out[:a] + out[j:]
    # Then any rule the cuts left with nothing in it, and the line it sat on.
    gone = []
    src = blanked(out)
    for sel, a, b, _ in rules(out):
        if src[a:b].strip() == "" and not sel.startswith("@"):
            gone.append((a - len(sel) - 1, b + 1))
    for a, b in reversed(gone):
        j = b
        while j < len(out) and out[j] == "\n":
            j += 1
        out = out[:a] + out[j:]
    return out, len(gone)


def main():
    what = sys.argv[1] if len(sys.argv) > 1 else "list"
    raw = open(STYLESHEET).read()
    if what == "list":
        found = dead(raw)
        for line, sel, prop, val, _, _ in found:
            print(f"line {line:5d}  {sel.split(',')[0][:30]:32s} {prop}: {val[:40]}")
        print(len(found), "dead declarations")
        return
    if what == "strip":
        cut = len(dead(raw))
        out, emptied = strip(raw)
        open(STYLESHEET, "w").write(out)
        print(cut, "declarations removed,", emptied, "rules left empty and removed")
        return
    if what == "same":
        before, after = resolved(open(sys.argv[2]).read()), resolved(raw)
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

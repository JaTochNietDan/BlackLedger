# Family headquarters and crew operations

User amendment, September14: a family needs an owned business as its base, the
player chooses that base when forming the family, and hired people should be
able to perform the practical work the player can perform from orders issued
there. This extends the active game goal; it does not replace earlier scope.

## Formation and bases

Implemented in this workstream: family formation is an explicit action at the
chosen business, requiring one owned trading business,25 respect, usable
premises, and independence from another family. Acquiring a second business or
waiting no longer silently incorporates the player. The business remains a
working concern with its existing costs and earnings. Establishing or moving
the base takes30 game minutes. An existing player family gets a suitable base
when upgrading a save; its later loss requires an explicit replacement choice.
NPC families keep a stable base rather than moving whenever they acquire a
richer property. Succession retains the family's headquarters deed, including
the player's surviving organization after death. NPC families choose a remaining
business if their base is lost. Losing the final business leaves no valid base;
it does not invent a deed or silently dissolve the family.

The base is an actual address used by the existing family home/rally function,
not the leader's residence. Housing remains separate. Public family information
shows a known family's base; personal homes are not exposed by this field.

## Current operations implementation

Saved named orders now cover restocking and assassination, with outbound/work/
return stages, parallel operatives, cash reservations, recall, arrival/commit
validation and HTTP retry coverage. The Business panel at headquarters and the
Families page contain an order selector and register. Assassination follows a
recent sighting rather than hidden movement and shares combat/scene effects.
Restocking shares the player supply effect. Losing the issuing life cancels
orders; surviving operatives retain unused budget. This is not yet the full
succession-continuation design below. Other operation adapters, comprehensive
resource custody and all remaining practical-command coverage are outstanding.

## Operations expansion to implement

The current code still contains old direct/first-associate delegation and a
single collections task. Those are not the requested operations system. Replace
those paths with an authoritative saved order per named operative; do not
pretend that exposing those buttons at headquarters completes the request.

An operations panel in the headquarters interior should let the player choose
one or more available hired people, an operation, target and budget/equipment.
Show expected outbound/work/return time, committed costs, prerequisites and
risk before dispatch. List active orders and results there and in People.
Target selection may inspect the public city directory while the player stays
at the base. Navigation and reading must not advance simulation time.

Use the union of signed family members and the original paid associate. Being
first in Player.Crew must not determine who acts or dies. A person cannot be on
two missions, guard a business and work elsewhere simultaneously, or receive an
order while dead, held, injured beyond the job's limit or away on another job.
A hired person's own weapon, skill, trust and injuries govern their attempt.
Use existing travel/arrival state and public journeys. Dispatch does not
teleport the player to the target or swap their identity into the NPC.

Order lifecycle: accepted/preparing, outbound, working, returning, settled (or
failed/cancelled). Save actor IDs, issuer family/life, target, stage, timestamps,
reserved money/items, and result. Reserve costs once and refund only unused
resources once. Recheck the target at arrival and resolution: dead/moved victim,
changed owner, ruined building, depleted cash or an unavailable operative can
invalidate the job. A failed attempt can injure, arrest or kill the operative;
the player can receive attention and retaliation without taking their wounds.
Do not clear the entire crew when one person dies. Resolved incidents must
carry the actual attacker and correct street/home/interior setting.

Queued work must survive reloads and exactly-once HTTP retries. Already-issued
orders do not become newly issued work when headquarters moves or is lost.
Recall acts on a real journey; it cannot undo a committed robbery or explosion.
On leader death, distinguish surviving-family succession from the next player
life; no money or authority should transfer to the new protagonist for free.

## Command coverage

| Work | Required crew support |
| --- | --- |
| Rob premises, mug, burgle, strip vehicles | Named operative, actual target funds/items, recognition and retaliation |
| Assassinate, including at home | Known target location, approach/travel, target availability, actual attacker and indoor cues |
| Bomb, incendiary attack, building drive-by, sabotage | Reserved charges/cash/ammunition/fuel, appropriate people/equipment, shared damage rules |
| Collections, deliveries, procurement, sale of goods | Cargo/cash custody, capacity, market changes, travel and funded settlement |
| Repair, restock, staffing and other business errands | Owner authorization, budget cap, actual work duration and shared business rules |
| Guarding, surveillance, intimidation and other practical person/place work | Explicit role/busy-state conflict, discovered public information and existing consequence rules |

Personal decisions remain personal: formation and succession choices, dialogue
commitments on the player's behalf, residence changes, clothing the player,
medical recovery, serving a sentence, new life, and playing the player's own
minigame. Crew can be sent for the practical errand around such a decision
(e.g. procurement); they do not become the protagonist or play a card hand in
their place. This is an actor/authority distinction, not a reason to leave
practical operations unimplemented.

Prefer reusable action handlers with an explicit actor/issuer/resources context
and shared readiness/effect rules. Avoid one generic proxy that overwrites
World.Player, and avoid copied crime implementations that drift from the
player's rules. First reconcile legacy Hand and Task assumptions, then wire the
operation catalogue and UI through normal Execute/API idempotency.

## Verification required before release

Explicit formation/selection/reload and both kinds of succession; lost/repaired
headquarters; multiple independent concurrent orders; actor death/arrest/recall;
target movement/ownership changes; no extra player teleport/damage/cash creation;
per-kind resource and consequence coverage; reload/retry/old-save coverage;
working operations UI at desktop and compact width; actual indoor incident
playback. Existing formation checks alone do not establish dispatch completion.

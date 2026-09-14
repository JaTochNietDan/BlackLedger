# Mariner lodging integration

Current requirement: objective903a4109-a113-4a4f-b726-f6ef277f7e41. The Mariner must be acquirable and manageable, with income from living NPC renters. This document records integration findings, not completed acceptance.

## Existing seams and constraints

- The Mariner is location `room`, type `home`, cost0. `move_home` uses this cost; changing it to a business purchase price would charge the freehold price when the player merely rents a room. Keep residential and business acquisition prices distinct.
- `Actions` opens business controls when Property.Income is positive. Acquisition, holdings, family income, forfeiture and operating obligations also inspect business eligibility; adding only a Buy button would leave several systems inconsistent.
- NPC.Post is a workplace, Location is current physical presence, and Heading is temporary travel. None is a lease. Create persisted tenant identity/lease state rather than charging everyone seen in the lobby or treating absent tenants as vacant rooms.
- `TradeOf` uses location Kind. Lodging needs staffing, supplies, remedy and management rules through this existing trade system, while retaining type home and player residential actions.
- The existing room daily charge is HomeRent, used by both DailyCost and Books. Owner occupancy must be handled explicitly to avoid paying rent to oneself.
- NPC.Purse and PayTheCity already represent money. PayTheCity in core/mugging.go currently bundles room and food in LivingCost=12, offset by daily earnings; split that existing room allowance rather than adding a second rent debit. PayDay in wages.go adjusts worker trust, not purse cash. Rent collection must debit renters and credit the holder once; dead tenants must neither pay nor produce income. It must not mint duplicate hourly takings on top of collected rent.
- SetOut/wanted currently direct people to posts, obligations and errands. A lease needs a return-home reason without overriding custody, urgent work or ongoing journeys.
- Books omitted RoomTrade and CollectionShare. A shared HourlyIncome now removes that discrepancy before adding a lodging rate; targeted tests pass, core/store/HTTP integration suites passed.

## Remaining implementation and acceptance

1. Persist leases with clear capacity, rate and payment/departure semantics; migrate old saves idempotently without stealing property or reviving/deleting NPCs.
2. Admit actual living renters, preserve their lease while at work/travelling, and release rooms on death/departure. Drive nighttime return travel from the lease.
3. Offer acquisition separately from residential rental and reuse staffing, supply, maintenance, arrears and manager controls.
4. Collect rent exactly once from tenants' money and show the occupancy-derived rate consistently in property details and daily accounts. Test zero occupancy, death, insufficient funds, vacancy/refill, ownership changes and player-owned residence costs.
5. Display named tenants, capacity, rate, collected income and problems in a styled lodging register. Keep all information server-authoritative.
6. Exercise acquisition and management through the HTTP command API on a fresh isolated fixture; verify saved/reloaded leases, clock partition invariance, obligations and browser display. Never QA against the main campaign.

No Mariner feature is enabled yet. Broad visual, interior, assassination and production requirements remain active alongside this work.

The latest amendment broadens this to city-wide homes appropriate to wealth and standing. Establish housing capacity and price tiers first, then lease the Mariner through the same assignment system. Employment/Post must remain independent of residence. Existing core/mugging.go contains unrelated working-tree edits; preserve them while coordinating any PayTheCity integration.

#!/usr/bin/env bash
set -euo pipefail
visual_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$visual_root"
visual_scenario="${1:-damage}"
case "$visual_scenario" in police|damage|warning|russo-warning|attack|voice) ;; *) echo 'Choose police, damage, warning, russo-warning, attack or voice.' >&2; exit 2;; esac
visual_go="$(command -v go || true)"
if [[ -z "$visual_go" && -x /usr/local/go/bin/go ]]; then visual_go=/usr/local/go/bin/go; fi
if [[ -z "$visual_go" ]]; then echo 'Go is required for the isolated preview.' >&2; exit 1; fi
if [[ ! -d node_modules ]]; then npm ci; fi
npm run build
mkdir -p .runtime
visual_db="$visual_root/.runtime/visual-preview-$visual_scenario.sqlite3"
if [[ ! -e "$visual_db" ]]; then "$visual_go" run ./cmd/qa-fixture "$visual_db" "$visual_scenario"; fi
"$visual_go" build -o .runtime/visual-preview-server ./cmd/blackledger
export BLACK_LEDGER_DB="$visual_db"
export BLACK_LEDGER_WEB="$visual_root/dist"
export BLACK_LEDGER_PORT="${BLACK_LEDGER_PORT:-8840}"
printf 'Visual preview: http://127.0.0.1:%s (%s fixture)\n' "$BLACK_LEDGER_PORT" "$visual_scenario"
exec .runtime/visual-preview-server

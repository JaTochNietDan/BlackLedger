#!/usr/bin/env bash
# Always create a fresh, isolated save. Never touch the campaign database.
set -euo pipefail
city_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$city_root"
city_scenario="${1:-city3d}"
case "$city_scenario" in city3d|city3d-night|city3d-blast|city3d-traffic|city3d-junction) ;; *) echo 'Choose city3d, city3d-night, city3d-blast, city3d-traffic, or city3d-junction.' >&2; exit 2;; esac
mkdir -p .runtime
city_save="$city_root/.runtime/$city_scenario-$(date +%Y%m%d-%H%M%S)-$$.sqlite3"
if [[ ! -d node_modules ]]; then npm ci; fi
npm run build
go run ./cmd/qa-fixture "$city_save" "$city_scenario"
go build -o .runtime/city3d-preview-server ./cmd/blackledger
export BLACK_LEDGER_WEB="$city_root/dist"
printf 'Isolated 3D city: http://127.0.0.1:%s\nSave: %s\n' "${BLACK_LEDGER_PORT:-8840}" "$city_save"
exec .runtime/city3d-preview-server -db "$city_save" -addr ":${BLACK_LEDGER_PORT:-8840}"

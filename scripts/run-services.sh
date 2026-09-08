#!/usr/bin/env bash
# Start the local AI and voice services Black Ledger talks to.
# Both are optional: authored fallbacks play without them.
set -euo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"
mkdir -p .runtime

export OLLAMA_HOST=127.0.0.1:11435
export OLLAMA_MODELS="$root/.tools/models"
export OLLAMA_NO_CLOUD=1

up() { curl -fsS --max-time 2 "$1" >/dev/null 2>&1; }

if up http://127.0.0.1:11435/api/tags; then
  echo "ollama: already running on 11435"
elif [[ -x .tools/ollama/ollama ]]; then
  nohup .tools/ollama/ollama serve >> .runtime/ollama.log 2>&1 &
  for _ in $(seq 1 40); do up http://127.0.0.1:11435/api/tags && break; sleep 0.5; done
  up http://127.0.0.1:11435/api/tags && echo "ollama: started on 11435" || echo "ollama: FAILED (see .runtime/ollama.log)" >&2
else
  echo "ollama: .tools/ollama/ollama missing" >&2
fi

if up http://127.0.0.1:8787/health; then
  echo "voice: already running on 8787"
else
  nohup python3 director/server.py >> .runtime/director.log 2>&1 &
  for _ in $(seq 1 40); do up http://127.0.0.1:8787/health && break; sleep 0.5; done
  up http://127.0.0.1:8787/health && echo "voice: started on 8787" || echo "voice: FAILED (see .runtime/director.log)" >&2
fi

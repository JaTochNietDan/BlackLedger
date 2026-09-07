# Black Ledger

Separate browser mafia prototype; does not modify Afterlight or its saves.

Run `python3 mafia-game/server.py` from the repository root, then open http://127.0.0.1:8791.

State: `mafia-game/.runtime/campaign.sqlite3`. Override with BLACK_LEDGER_DB for isolated tests. Bind is localhost only. AI uses BLACK_LEDGER_OLLAMA (default http://127.0.0.1:11435), BLACK_LEDGER_MODEL (default qwen3:14b); voice proxies the existing local service at AFTERLIGHT_DIRECTOR_URL (default http://127.0.0.1:8787). Neither service is required to play authored scenarios.

Run tests: `python3 -m unittest discover -s mafia-game/tests -v`.

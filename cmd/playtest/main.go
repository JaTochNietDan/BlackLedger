// The workshop server.
//
// A separate binary on a separate port, deliberately. The inbox asked for a
// debug mode "separate from the main game", and the strongest form of separate
// is that the game's binary contains none of this: there is no flag on the game
// that turns it on, no route to find, and nothing to leave switched on by
// accident in something a player runs.
//
// It serves a scratch world it builds for itself. It will not open a campaign,
// and refuses to start if it is pointed at one — the whole value of a workshop
// is that you can do anything in it, which is only true if nothing in it
// matters.
//
//	mise run workshop
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"blackledger/core"
)

func main() {
	port := flag.String("port", envOr("BLACK_LEDGER_WORKSHOP_PORT", "8790"), "port to listen on")
	seed := flag.Int("seed", 4242, "which scratch world to build")
	flag.Parse()

	// A workshop that could open a campaign is a workshop that can ruin one.
	if db := os.Getenv("BLACK_LEDGER_DB"); db != "" {
		log.Fatalf("the workshop does not open saves, and BLACK_LEDGER_DB is set to %q.\n"+
			"It builds a scratch world of its own so that nothing done in it matters.", db)
	}

	world := core.New(uint32(*seed))
	root := webRoot()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		// Rebuilt per request from the same seed, so the workshop cannot drift
		// into a state nobody can reproduce.
		send(w, core.New(uint32(*seed)).Public())
	})
	mux.HandleFunc("/api/kinds", func(w http.ResponseWriter, r *http.Request) {
		send(w, kinds())
	})
	mux.Handle("/", page(root))

	log.Printf("workshop on http://127.0.0.1:%s — scratch world %d, %d addresses, no save",
		*port, *seed, len(world.Public()["locations"].([]map[string]any)))
	log.Fatal(http.ListenAndServe("127.0.0.1:"+*port, mux))
}

// kinds is every moment the city knows how to show, ordered by what it is worth
// stopping for. Read from the core rather than listed here: a kind added to the
// game turns up in the workshop without anybody remembering to add it.
func kinds() []map[string]any {
	out := []map[string]any{}
	for _, k := range core.MomentKinds() {
		out = append(out, map[string]any{
			"kind": k, "gravity": core.Gravity(k), "hold": core.Hold(k),
		})
	}
	return out
}

// page serves the built workshop entry. Everything else falls through to the
// same assets the game uses, because the workshop is the game's components.
func page(root string) http.Handler {
	files := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := strings.TrimPrefix(r.URL.Path, "/")
		if clean == "" || clean == "index.html" {
			http.ServeFile(w, r, filepath.Join(root, "playtest.html"))
			return
		}
		if _, err := os.Stat(filepath.Join(root, clean)); err != nil {
			http.ServeFile(w, r, filepath.Join(root, "playtest.html"))
			return
		}
		files.ServeHTTP(w, r)
	})
}

func webRoot() string {
	if root := os.Getenv("BLACK_LEDGER_WEB"); root != "" {
		return root
	}
	if _, err := os.Stat(filepath.Join("dist", "playtest.html")); err == nil {
		return "dist"
	}
	fmt.Fprintln(os.Stderr, "no built workshop page found; run `npm run build` first")
	return "dist"
}

func send(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

// The map the city is dressed with.
//
// Where each address stands is not in here and never will be: that comes from
// the core's own coordinates, and a second file allowed to disagree about it
// would be a second city. What is in here is everything the core has no opinion
// about — which picture stands in which slot, and how far it is nudged — so it
// can be arranged by hand without any of it being able to contradict the world.
//
// It is a plain file in the repository on purpose. One person drags a building
// in the browser, the other edits the same JSON in an editor, and the change
// arrives as an ordinary diff either way.

const layoutFile = "art/city-layout.json"

// editable reports whether this build may be edited. Writing to the repository
// from a web request is a development convenience and has no business being
// reachable from a running game, so it is off unless it is asked for.
func editable() bool { return os.Getenv("BLACK_LEDGER_EDIT") == "1" }

func readLayout(w http.ResponseWriter) {
	body, err := os.ReadFile(layoutFile)
	if os.IsNotExist(err) {
		reply(w, 200, map[string]any{"slots": map[string]any{}, "editable": editable()})
		return
	}
	if err != nil {
		fail(w, 500, err)
		return
	}
	var saved map[string]any
	if err := json.Unmarshal(body, &saved); err != nil {
		fail(w, 500, err)
		return
	}
	saved["editable"] = editable()
	reply(w, 200, saved)
}

func writeLayout(w http.ResponseWriter, r *http.Request) {
	if !editable() {
		fail(w, 403, errNotEditable)
		return
	}
	var incoming struct {
		Slots map[string]struct {
			Sprite string  `json:"sprite"`
			DX     float64 `json:"dx,omitempty"`
			DY     float64 `json:"dy,omitempty"`
		} `json:"slots"`
	}
	if err := body(r, &incoming); err != nil {
		fail(w, 400, err)
		return
	}
	if incoming.Slots == nil {
		incoming.Slots = map[string]struct {
			Sprite string  `json:"sprite"`
			DX     float64 `json:"dx,omitempty"`
			DY     float64 `json:"dy,omitempty"`
		}{}
	}
	out, err := json.MarshalIndent(map[string]any{"slots": incoming.Slots}, "", " ")
	if err != nil {
		fail(w, 500, err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(layoutFile), 0o755); err != nil {
		fail(w, 500, err)
		return
	}
	// Written whole rather than appended to, so a half-finished save cannot
	// leave a file that parses but describes a city nobody arranged.
	if err := os.WriteFile(layoutFile, append(out, '\n'), 0o644); err != nil {
		fail(w, 500, err)
		return
	}
	reply(w, 200, map[string]any{"saved": len(incoming.Slots)})
}

type layoutError string

func (e layoutError) Error() string { return string(e) }

const errNotEditable = layoutError(
	"this build is not editable; start it with BLACK_LEDGER_EDIT=1 to arrange the map")

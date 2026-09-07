package main

import (
	"blackledger/core"
	"fmt"
	"strings"
)

func repeatedProposal(w *core.World, p core.Proposal) error {
	normalize := func(s string) string {
		return strings.ToLower(strings.Join(strings.Fields(strings.Trim(s, " \n\t\"“”")), " "))
	}
	// Current-life memory includes declined/abandoned offers: those should not be re-pitched verbatim either.
	for _, m := range w.Arrangements {
		if m.Life != w.Life {
			continue
		}
		if normalize(m.Title) == normalize(p.Title) || normalize(m.Offer) == normalize(p.Body) {
			return fmt.Errorf("proposal repeats a recent offer %q; write a distinct task and title, preserving any required contact connection", m.Title)
		}
	}
	return nil
}

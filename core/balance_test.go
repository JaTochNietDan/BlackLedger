package core

import "testing"

// The balance suites are the slow half of this project by a wide margin: 54
// files of them run tens of thousands of simulated days, and measured with
// -json they account for some 676 seconds of test time against a whole suite
// that finishes in 115 because they run in parallel with each other.
//
// That is the right cost to pay before a commit and the wrong one to pay while
// iterating on a single change. `heavy` skips them under -short, so
// `go test -short ./core` answers "did I break anything structural" in seconds
// and the full run still answers "did I move the balance" before anything is
// committed. Nothing is skipped in the gate, which never passes -short.
func heavy(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("balance: run without -short to measure it")
	}
	t.Parallel()
}

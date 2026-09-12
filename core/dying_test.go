package core

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// Eleven rules in this game end a life and nothing could count which of them
// fired. The cause each one writes is a sentence with a name and an address in
// it — "An attempt on Tilda Gruber at The Mariner that went the other way" —
// which reads properly in a ledger and cannot be grouped: two hundred and sixty
// campaigns reported eleven distinct causes where there were three.
//
// So a death carries the rule as well as the sentence, and this makes sure the
// twelfth way to die carries it too. `Die` still exists, because a save written
// before this did not record one, and it is not for new code.
func TestEveryWayToDieNamesTheRuleThatDidIt(t *testing.T) {
	t.Parallel()
	found, unnamed := 0, []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for i, line := range strings.Split(string(body), "\n") {
			switch {
			case strings.Contains(line, "w.DieOf("):
				found++
			case strings.Contains(line, "w.Die("):
				unnamed = append(unnamed, filepath.Base(path)+":"+strconv.Itoa(i+1)+" "+strings.TrimSpace(line))
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(unnamed) > 0 {
		t.Errorf("%d way(s) to die that cannot be counted — use DieOf and name the rule:\n  %s",
			len(unnamed), strings.Join(unnamed, "\n  "))
	}
	if found < 10 {
		t.Fatalf("only %d ways to die were found, so this is reading the wrong thing", found)
	}
	t.Logf("%d ways to die, every one of them named", found)
}

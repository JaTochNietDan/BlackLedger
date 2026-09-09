package core

import (
	"fmt"
	"strings"
)

// Organizations are named by three separate systems — the seeded families, the
// splinters, and whatever the player's own people end up being called — and the
// paper has to put all of them into sentences. Two things go wrong if nobody
// thinks about it, and both were printed before an issue was read end to end: a
// name that carries its own article begins a sentence in lower case ("the Rizzo
// Crew has ceased to operate"), and a name that is grammatically plural takes a
// singular verb ("Cesare Ferro's people has had a poor few weeks").

// Leads returns a name fit to begin a sentence, capitalising an article the
// name carries with it. Mid-sentence the name is correct as it stands.
func Leads(name string) string {
	if strings.HasPrefix(name, "the ") {
		return "The " + name[len("the "):]
	}
	return name
}

// pluralNames are the endings that make an organization take a plural verb.
var pluralNames = []string{"people", "Brothers", "Boys"} // not prose: organization name endings

// Agree picks the verb form that goes with an organization's name.
func Agree(name, singular, plural string) string {
	for _, ending := range pluralNames {
		if strings.HasSuffix(name, ending) {
			return plural
		}
	}
	return singular
}

// titled is true for a role that already carries its own complement — "head of
// the Russo Outfit", "runs Bluebird Laundry". A unique office named with what
// it is an office of takes no article, and adding one printed "They were a head
// of the Russo Outfit". This was introduced by the fix for "They were
// Lieutenant" and found by reading the paper end to end a second time.
func titled(role string) bool {
	return strings.Contains(role, " of ") || strings.HasPrefix(strings.ToLower(role), "runs ")
}

// article picks "a" or "an" for a word. The obituaries read "They were
// Lieutenant" until somebody read one.
func article(word string) string {
	if word == "" {
		return "a"
	}
	if strings.ContainsRune("aeiou", rune(strings.ToLower(word)[0])) {
		return "an"
	}
	return "a"
}

// upper1 begins a sentence with a phrase written for the middle of one.
func upper1(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// spelled writes a small number the way a newspaper does. Prose that says "2
// people in the same organization stood below them" reads like a report from a
// machine, which in this case it was. Above twelve the paper uses figures, as
// papers do.
var numbers = []string{"no", "one", "two", "three", "four", "five", "six",
	"seven", "eight", "nine", "ten", "eleven", "twelve"}

func spelled(n int) string {
	if n >= 0 && n < len(numbers) {
		return numbers[n]
	}
	return fmt.Sprintf("%d", n)
}

// withArticle is a role fit to follow "They were": an ordinary job takes an
// article, a titled office does not.
func withArticle(job string) string {
	if titled(job) {
		return job
	}
	return article(job) + " " + job
}

// counted agrees a noun with the number in front of it. The city was printing
// "$164 due in 1 days" on three buttons at the casino and "1 stories in today's
// paper are about you" at the Herald. Everything else here is keyed to code
// rather than to save state, and so is this: nothing about a count is stored.
func counted(n int, singular, plural string) string {
	if n == 1 {
		return "1 " + singular
	}
	return fmt.Sprintf("%d %s", n, plural)
}

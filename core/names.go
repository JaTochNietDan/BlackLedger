package core

import "strings"

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
var pluralNames = []string{"people", "Brothers", "Boys"}

// Agree picks the verb form that goes with an organization's name.
func Agree(name, singular, plural string) string {
	for _, ending := range pluralNames {
		if strings.HasSuffix(name, ending) {
			return plural
		}
	}
	return singular
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

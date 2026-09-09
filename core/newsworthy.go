package core

// How big a story is, as a newspaper would rank it.
//
// This is not the same question as how hard the police look, and the two were
// the same map. `scrutinyWeight` decides what a story costs the player in
// attention: a killing raises the city's temperature, a column about the price
// of coal does not. It was also being used to decide which story leads the
// paper and which one the city screen announces — and it does not cover half
// the kinds the paper actually files. An arrest, an attempt on somebody's life,
// a family splitting and a family moving back onto premises all scored zero,
// which is the same score as the weather.
//
// So an arrest could never lead an edition, and once the city page started
// filing filler every morning the banner on the city screen would announce the
// weather over a man being shot, because it simply took the newest story.
//
// Kept separate on purpose. Changing what a story costs in police attention
// changes the balance of the game; changing what leads the paper changes only
// what the player is told first.

var newsworthiness = map[string]int{
	"killing":  100,
	"attempt":  85,   // somebody survived it, which is the only difference
	"war":      80,
	"attack":   70,
	"collapse": 65,   // an organization ending is the end of something
	"split":    60,
	"arrest":   55,
	"police":   50,
	"seizure":  45,
	"recovery": 35,
	"robbery":  30,
	"politics": 25,
	"business": 15,
	"obituary": 10,   // set apart from the news anyway, never a lead
	"civic":    5,    // the city page: real, and never the front page
}

// Newsworthiness is how far up the page a kind of story belongs. Anything the
// paper can file has a place; an unknown kind sits above the filler and below
// everything that happened to a person, which is where a new kind of story
// most likely belongs until somebody decides otherwise.
func Newsworthiness(kind string) int {
	if n, ok := newsworthiness[kind]; ok {
		return n
	}
	return 20
}

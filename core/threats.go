package core

import "fmt"

// What people say to your face.
//
// From the inbox: "I don't see a lot of cursing from characters in this game, we
// should increase that since it's with the mafia style. Characters should be
// able to make threats to you too, I have not seen that yet."
//
// The city already knew who had a reason to dislike the player and how badly —
// n.Sore, and n.SoreAt for what it was about — and did nothing with it but sort
// a list. Somebody carrying something against you and standing in the same room
// as you ought to say so, and the closer they are to acting on it the less
// polite the saying gets.
//
// Nothing here invents a grievance. Every line is built out of the weight the
// simulation already committed and the words the player has already earned.

const (
	// Muttering is the weight below which it is a look and a remark.
	Muttering = 25
	// Plain is where it stops being a remark and becomes a warning.
	Plain = 45
)

// oaths are period and mild by the standards of the men saying them. They are
// picked by the speaker's name so that the same man swears the same way twice.
var oaths = []string{"the hell", "damn", "christ's sake", "goddamn"}

func oath(id string) string { return oaths[int(FaceOf(id))%len(oaths)] }

// Threat is what somebody says to the player, or "" when they have nothing to
// say. Only from people the player actually knows: a stranger with a grievance
// is a stranger, and this game does not put words in the mouths of people it
// has not introduced.
func (w *World) Threat(n *NPC) string {
	if n == nil || n.Dead || n.Sore <= 0 || !w.Known(n) {
		return ""
	}
	about := n.SoreAt
	if about == "" {
		about = "what you did"
	}
	swear := oath(n.ID)
	hot := TemperamentOf(n)

	switch {
	case n.Sore >= GrudgeActs:
		// Past the weight at which people move. This is not a warning any more.
		if hot.ID == "careful" {
			return fmt.Sprintf("“I am not going to make a speech about %s. You will know when it is settled.”", about)
		}
		if hot.ID == "vain" {
			return fmt.Sprintf("“Everybody in this room knows about %s. They will know how it ended, too.”", about)
		}
		return fmt.Sprintf("“You had better hope somebody is watching your door, %s. I have not forgotten %s.”", swear, about)
	case n.Sore >= Plain:
		if hot.ID == "careful" {
			return fmt.Sprintf("“%s is not finished. I am telling you that once.”", capitalise(about))
		}
		return fmt.Sprintf("“Get out of my sight before I do something about %s, for %s.”", about, swear)
	case n.Sore >= Muttering:
		if hot.ID == "loyal" {
			return fmt.Sprintf("“My people know about %s. That is all I will say in here.”", about)
		}
		return fmt.Sprintf("“%s. You have some nerve standing there.”", capitalise(about))
	}
	return fmt.Sprintf("“%s. I have not said anything, and I am not going to.”", capitalise(about))
}

func capitalise(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	if r[0] >= 'a' && r[0] <= 'z' {
		r[0] = r[0] - 32
	}
	return string(r)
}

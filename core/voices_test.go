package core

import "testing"

// "Detective Harlow has a female voice but it's a male picture." The face came
// off a hash of an id in the interface and the voice came off a hash of a name
// in the speech service, and neither knew the other existed.

func TestAVoiceMatchesThePaintingItComesOutOf(t *testing.T) {
	women, men := 0, 0
	for _, l := range Locations {
		_ = l
	}
	w := New(23)
	for i := range w.NPCs {
		n := &w.NPCs[i]
		profile := w.VoiceOf(n.ID)
		at := voiceIndex(t, profile)
		looks := castLooks[FaceOf(n.ID)]
		if painted, ok := paintedLooks[n.ID]; ok {
			looks = painted
		}
		switch looks {
		case looksWoman:
			women++
			if at < womenFrom || at > womenTo {
				t.Errorf("%s is painted as a woman and speaks with %s", n.Name, profile)
			}
		default:
			men++
			if at < menFrom || at > menTo {
				t.Errorf("%s is painted as a man and speaks with %s", n.Name, profile)
			}
		}
	}
	t.Logf("%d people painted as women, %d as men, every one of them speaking to match", women, men)
	if women == 0 || men == 0 {
		t.Error("this city is all one thing")
	}
}

// The one the report was about.
func TestHarlowSoundsLikeTheManInTheTrenchCoat(t *testing.T) {
	w := New(23)
	if w.NPC("harlow") == nil {
		t.Skip("no detective in this city")
	}
	at := voiceIndex(t, w.VoiceOf("harlow"))
	if at < menFrom || at > menTo {
		t.Errorf("the detective is painted in a trench coat and speaks with cast2-%02d", at)
	}
}

// And a voice never moves. Somebody who sounded one way yesterday sounds that
// way for the rest of their life, whatever else happens to them.
func TestAVoiceIsFixedForLife(t *testing.T) {
	w := New(23)
	n := &w.NPCs[3]
	was := w.VoiceOf(n.ID)
	n.Name, n.Role, n.Faction, n.Rank = "Somebody Else", "Boss", "russo", RankLeader
	if now := w.VoiceOf(n.ID); now != was {
		t.Errorf("they were %s and are now %s", was, now)
	}
}

func voiceIndex(t *testing.T, profile string) int {
	t.Helper()
	if len(profile) != 8 || profile[:6] != "cast2-" {
		t.Fatalf("%q is not a voice the service holds", profile)
	}
	at := int(profile[6]-'0')*10 + int(profile[7]-'0')
	if at < 0 || at > 47 {
		t.Fatalf("%q is outside the cast", profile)
	}
	return at
}

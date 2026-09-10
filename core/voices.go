package core

// Whose face somebody has, and what they sound like — decided in one place,
// because they are the same decision. The city drew faces from a hash of an id
// and the voice service drew voices from a hash of a name, and the two hashes
// knew nothing about each other: Detective Harlow is painted in a trench coat
// with a five o'clock shadow and spoke in a woman's voice for as long as the
// game has had voices at all.
//
// Nobody in Bellwether is gendered by the rules and the prose still says "they"
// about everybody. This is not about the rules. It is about a painting and a
// recording being of the same person.

// The paintings. Six were made by hand and the rest come off a sheet of
// twenty-four; what is written here is what is in the picture, read off the
// canvas rather than guessed from the prompt that made it.
const (
	looksWoman = 'f'
	looksMan   = 'm'
)

var castLooks = [CastFaces]byte{
	looksMan, looksMan, looksWoman, looksMan, looksWoman, looksMan,
	looksWoman, looksMan, looksWoman, looksMan, looksMan, looksWoman,
	looksMan, looksWoman, looksMan, looksMan, looksWoman, looksMan,
	looksMan, looksMan, looksWoman, looksMan, looksWoman, looksWoman,
}

// The six painted by hand, by the id the interface draws them for.
var paintedLooks = map[string]byte{
	"mara": looksWoman, "leo": looksMan, "vittorio": looksMan,
	"elena": looksWoman, "harlow": looksMan,
}

// FaceOf is which of the sheet's faces somebody wears. The interface used to
// work this out for itself; it is the core's now, because the voice has to
// agree with it and only one of them can be the author of it.
//
// The arithmetic is FNV-1a over the id, which is what the interface was already
// doing — this is the same answer for every person in every save, moved to the
// side of the wall that owns facts.
func FaceOf(id string) int {
	h := uint32(2166136261)
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= 16777619
	}
	return int(h % uint32(CastFaces))
}

// The voice profiles the speech service holds, in the bands it built them in:
// twenty American women, then fourteen American men. British voices exist and
// are left for anybody who should sound like an import.
const (
	womenFrom, womenTo = 0, 19
	menFrom, menTo     = 20, 33
)

// VoiceProfile is which recorded voice somebody speaks in. It follows the
// picture: a face painted as a woman is given a woman's voice and a face
// painted as a man is given a man's, and which one inside that band is the
// person's own business, fixed for life by their id.
func VoiceProfile(id string, face int) string {
	looks := byte(0)
	if painted, ok := paintedLooks[id]; ok {
		looks = painted
	} else if face >= 0 && face < CastFaces {
		looks = castLooks[face]
	}
	from, to := menFrom, menTo
	if looks == looksWoman {
		from, to = womenFrom, womenTo
	}
	h := uint32(2166136261)
	for i := 0; i < len(id); i++ {
		h ^= uint32(id[i])
		h *= 16777619
	}
	at := from + int(h%uint32(to-from+1))
	return "cast2-" + string(rune('0'+at/10)) + string(rune('0'+at%10))
}

// VoiceOf is the voice of somebody in this city, by their id.
func (w *World) VoiceOf(id string) string {
	// The player's own face is theirs to choose, so their voice follows the
	// choice rather than the hash.
	if id == w.Player.Name && w.Player.Face > 0 {
		return VoiceProfile(id, (w.Player.Face-1)%CastFaces)
	}
	return VoiceProfile(id, FaceOf(id))
}

package core

import "testing"

// The director may set a fact better. It may not add one.

const plainWeather = "Rain through the day and into the evening. The gutters on the lower streets are carrying more than they were built for."

func TestAPolishThatAddsANumberIsRefused(t *testing.T) {
	// The dangerous invention: a reader takes a number as fact and acts on it.
	body, took := AcceptPolish(plainWeather,
		"Rain fell for 11 hours across the district, and 4 streets were closed by standing water near the harbour.")
	if took {
		t.Error("a rewrite invented a count of hours and streets and was accepted")
	}
	if body != plainWeather {
		t.Error("a refused rewrite did not fall back to what the facts said")
	}
}

func TestAPolishThatAddsANameIsRefused(t *testing.T) {
	if _, took := AcceptPolish(plainWeather,
		"Rain through the day. Commissioner Vance said the drains would be looked at when there was money for it."); took {
		t.Error("a rewrite put a named official into a paragraph about the weather")
	}
}

func TestAPolishThatAddressesTheReaderIsRefused(t *testing.T) {
	if _, took := AcceptPolish(plainWeather,
		"Rain through the day and into the evening. If you are going out tonight, take a coat and watch the low end of the harbour road."); took {
		t.Error("the paper started talking to the reader")
	}
}

func TestAGoodPolishIsTaken(t *testing.T) {
	better := "Rain from first light and no let up by dark. On the lower streets the gutters took more than they were built to take, and stood in the road for it."
	body, took := AcceptPolish(plainWeather, better)
	if !took {
		t.Fatal("a rewrite that adds nothing and reads better was refused")
	}
	if body != better {
		t.Error("the accepted rewrite was not the text returned")
	}
}

func TestTheModelExplainingItselfIsStripped(t *testing.T) {
	body, took := AcceptPolish(plainWeather,
		"Here is the rewritten paragraph in the requested register:\n\nRain from first light and no let up by dark, and the gutters on the lower streets took more than they were built to take.")
	if !took {
		t.Fatal("the copy was thrown away with the preamble")
	}
	if len(body) > 200 || body[0] != 'R' {
		t.Errorf("the preamble survived: %q", body)
	}
}

func TestALongRambleIsRefused(t *testing.T) {
	long := ""
	for i := 0; i < 40; i++ {
		long += "Rain fell and the gutters carried it away again into the harbour. "
	}
	if _, took := AcceptPolish(plainWeather, long); took {
		t.Error("a rewrite of eight paragraphs was accepted for a two sentence brief")
	}
}

func TestOnlyTodaysUnpolishedBriefIsOffered(t *testing.T) {
	w := New(81)
	w.Minute = 5 * 1440
	w.CityPageDay()
	first, ok := w.NextPolish()
	if !ok {
		t.Fatal("today's page had nothing for the director to look at")
	}
	if !w.SetPolish(first.ID, "Rewritten copy that adds nothing at all.", true) {
		t.Fatal("the rewrite could not be written back")
	}
	second, ok := w.NextPolish()
	if !ok || second.ID == first.ID {
		t.Error("the same brief was offered twice")
	}
	// Yesterday's paper has gone out and is not reopened.
	w.Minute += 1440
	if _, ok := w.NextPolish(); ok {
		t.Error("the director was offered a brief from an issue already printed")
	}
}

func TestARefusalIsFinal(t *testing.T) {
	w := New(83)
	w.CityPageDay()
	s, _ := w.NextPolish()
	w.SetPolish(s.ID, "", false)          // refused
	for _, n := range w.News {
		if n.ID == s.ID {
			if !n.Polished {
				t.Error("a refused brief would be sent to the model again forever")
			}
			if n.Body == "" {
				t.Error("a refusal emptied the brief instead of leaving the facts")
			}
		}
	}
}

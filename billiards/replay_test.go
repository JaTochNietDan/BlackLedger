package billiards

import (
	"math"
	"strings"
	"testing"
)

func TestCompactReplayKeepsEveryImpactAndPose(t *testing.T) {
	result, err := New().Shoot(Rack(), Shot{Angle: math.Pi / 2, Speed: 7})
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := EncodeReplay(result)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeReplay(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded.Frames) != len(result.Frames) || len(decoded.Events) != len(result.Events) {
		t.Fatal("replay discarded impact data")
	}
	for i, f := range result.Frames {
		got := decoded.Frames[i]
		near(t, got.Time, f.Time, 1e-12)
		for j, b := range f.Balls {
			pose := got.Balls[j]
			if int(pose[0]) != b.ID {
				t.Fatal("changed ball identity")
			}
			near(t, float64(pose[1]), b.Position.X, 1e-6)
			near(t, float64(pose[2]), b.Position.Y, 1e-6)
			near(t, float64(pose[3]), b.Position.Z, 1e-6)
			for k := 0; k < 4; k++ {
				near(t, float64(pose[4+k]), b.Orientation[k], 1e-6)
			}
			pocket := 0
			if b.Pocketed {
				pocket = b.Pocket + 1
			}
			if int(pose[8]) != pocket {
				t.Fatal("changed pocket state")
			}
		}
	}
	if len(encoded) > 100000 {
		t.Fatal("break replay is too large for receipts", len(encoded))
	}
	t.Logf("break replay: %d encoded bytes, %d frames, %d events", len(encoded), len(decoded.Frames), len(decoded.Events))
}
func TestReplayDecoderRejectsInvalidAndOversizePayloads(t *testing.T) {
	for _, input := range []string{"", "not base64", "aGVsbG8=", strings.Repeat("a", maxEncodedReplayBytes+1)} {
		if _, err := DecodeReplay(input); err == nil {
			t.Fatal("accepted invalid replay")
		}
	}
}

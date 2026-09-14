package billiards

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// Replay omits velocities and repeated object keys. Poses contain
// [id,x,y,z,qx,qy,qz,qw,pocket], with pocket=0 for live balls and 1..6 otherwise.
// Float32 position accuracy is finer than the solver's convergence tolerance.
// Time remains float64 so simultaneous/nearby impacts retain their ordering.
type Replay struct {
	Version  int           `json:"v"`
	Duration float64       `json:"duration"`
	Frames   []ReplayFrame `json:"frames"`
	Events   []Event       `json:"events"`
}
type ReplayFrame struct {
	Time  float64      `json:"t"`
	Balls [][9]float32 `json:"balls"`
}

const maxReplayBytes = 8 << 20
const maxEncodedReplayBytes = 2 << 20

func EncodeReplay(r Result) (string, error) {
	tape := Replay{Version: 1, Duration: r.Duration, Events: r.Events}
	for _, f := range r.Frames {
		frame := ReplayFrame{Time: f.Time}
		for _, b := range f.Balls {
			pocket := 0
			if b.Pocketed {
				pocket = b.Pocket + 1
			}
			q := b.Orientation
			frame.Balls = append(frame.Balls, [9]float32{float32(b.ID), float32(b.Position.X), float32(b.Position.Y), float32(b.Position.Z), float32(q[0]), float32(q[1]), float32(q[2]), float32(q[3]), float32(pocket)})
		}
		tape.Frames = append(tape.Frames, frame)
	}
	raw, err := json.Marshal(tape)
	if err != nil {
		return "", err
	}
	if len(raw) > maxReplayBytes {
		return "", errors.New("shot replay is too large")
	}
	var out bytes.Buffer
	writer := zlib.NewWriter(&out)
	if _, err = writer.Write(raw); err != nil {
		return "", err
	}
	if err = writer.Close(); err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(out.Bytes())
	if len(encoded) > maxEncodedReplayBytes {
		return "", errors.New("compressed shot replay is too large")
	}
	return encoded, nil
}

// DecodeReplay is for verification and tools. HTTP clients may decode the
// base64/zlib JSON themselves; the server never accepts replay as a shot result.
func DecodeReplay(encoded string) (Replay, error) {
	var tape Replay
	if len(encoded) > maxEncodedReplayBytes {
		return tape, errors.New("encoded replay is too large")
	}
	packed, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return tape, err
	}
	reader, err := zlib.NewReader(bytes.NewReader(packed))
	if err != nil {
		return tape, err
	}
	defer reader.Close()
	raw, err := io.ReadAll(io.LimitReader(reader, maxReplayBytes+1))
	if err != nil {
		return tape, err
	}
	if len(raw) > maxReplayBytes {
		return tape, errors.New("expanded replay is too large")
	}
	if err = json.Unmarshal(raw, &tape); err != nil {
		return tape, err
	}
	if tape.Version != 1 {
		return tape, errors.New("unsupported replay version")
	}
	return tape, nil
}

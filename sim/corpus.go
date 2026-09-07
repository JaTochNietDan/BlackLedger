package sim

import (
	"blackledger/core"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

// ReadCorpus accepts data only. No recorded text is interpreted as an instruction.
func ReadCorpus(reader io.Reader) ([]core.Proposal, string, error) {
	const maxBytes = 2 << 20
	b, err := io.ReadAll(io.LimitReader(reader, maxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if len(b) > maxBytes {
		return nil, "", fmt.Errorf("corpus exceeds 2 MiB")
	}
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields()
	var proposals []core.Proposal
	if err = dec.Decode(&proposals); err != nil {
		return nil, "", err
	}
	var extra any
	if err = dec.Decode(&extra); err != io.EOF {
		return nil, "", fmt.Errorf("corpus must contain one JSON array")
	}
	if len(proposals) == 0 || len(proposals) > 128 {
		return nil, "", fmt.Errorf("corpus must contain 1–128 proposals")
	}
	w := core.New(27)
	for i, p := range proposals {
		if _, err = w.ValidateProposal(p); err != nil {
			return nil, "", fmt.Errorf("proposal %d: %w", i+1, err)
		}
	}
	digest := sha256.Sum256(b)
	return proposals, hex.EncodeToString(digest[:]), nil
}

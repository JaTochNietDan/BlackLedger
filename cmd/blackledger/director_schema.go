package main

import "blackledger/core"

// Decoding constraints assist the local model; core validation remains mandatory.
func proposalSchema(w *core.World, operation string, connection *core.ArrangementMemory) map[string]any {
	text := func(min, max int) map[string]any {
		return map[string]any{"type": "string", "minLength": min, "maxLength": max}
	}
	enum := func(values []string) map[string]any { return map[string]any{"type": "string", "enum": values} }
	speakers := []string{}
	beneficiaries := []string{""}
	for _, n := range w.NPCs {
		speakers = append(speakers, n.ID)
	}
	for _, f := range w.Factions {
		beneficiaries = append(beneficiaries, f.ID)
	}
	if connection != nil {
		speakers = []string{connection.Speaker}
		beneficiaries = []string{connection.Beneficiary}
	}
	return map[string]any{"type": "object", "additionalProperties": false,
		"required": []string{"title", "body", "speaker", "operation", "outcome", "beneficiary", "approaches"},
		"properties": map[string]any{
			"title": text(3, 70), "body": text(3, 1200), "speaker": enum(speakers), "operation": enum([]string{operation}), "outcome": text(3, 700), "beneficiary": enum(beneficiaries),
			"approaches": map[string]any{"type": "array", "minItems": 1, "maxItems": 2, "items": map[string]any{"type": "object", "additionalProperties": false, "required": []string{"method", "label"}, "properties": map[string]any{"method": enum([]string{"careful", "press"}), "label": text(3, 45)}}},
		}}
}

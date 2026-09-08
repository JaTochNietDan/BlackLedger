package core

import (
	"fmt"
	"strings"
)

// The director authors contextual labels, but cannot set costs or outcomes.
func addApproaches(scene *Scene, proposals []Approach) error {
	if len(proposals) > 2 {
		return fmt.Errorf("too many proposed approaches")
	}
	scene.Alternatives = map[string]Effect{}
	for _, approach := range proposals {
		label := strings.TrimSpace(approach.Label)
		if len(label) < 3 || len(label) > 65 {
			return fmt.Errorf("invalid approach label")
		}
		id := "approach:" + approach.Method
		if _, exists := scene.Alternatives[id]; exists {
			return fmt.Errorf("duplicate approach")
		}
		fx := scene.Effect
		switch approach.Method {
		case "careful":
			fx.Minutes += 30
			fx.Reward = max(0, fx.Reward-15)
			fx.Heat = max(0, fx.Heat-3)
		case "press":
			fx.Minutes = max(15, fx.Minutes-15)
			fx.Reward += 20
			fx.Heat += 5
		default:
			return fmt.Errorf("unsupported approach")
		}
		scene.Alternatives[id] = fx
		choice := Choice{ID: id, Label: label}.terms(fx)
		// Keep refusal last, after every actionable way to handle the job.
		last := len(scene.Choices) - 1
		scene.Choices = append(scene.Choices[:last], choice, scene.Choices[last])
	}
	return nil
}

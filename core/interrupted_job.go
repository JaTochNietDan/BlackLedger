package core

import "fmt"

type SuspendedJob struct {
	Scene     *Scene `json:"scene"`
	Remaining int    `json:"remaining"`
}

func (w *World) RunArrangement(job *Scene, remaining int) {
	w.RememberArrangement(job, "in_progress")
	start := w.Minute
	w.Advance(remaining)
	if w.Player.Alive && w.Event == nil {
		if w.Player.Heat+job.Effect.Heat >= 15 {
			w.PoliceStop(job)
		} else {
			w.CompleteArrangement(job)
		}
	} else if w.Player.Alive && w.Event != nil && w.Event.Kind == "business_pressure" {
		w.SuspendedJob = &SuspendedJob{Scene: job, Remaining: max(0, remaining-(w.Minute-start))}
		w.RememberArrangement(job, "paused")
		w.Log("Work put on hold", fmt.Sprintf("%s is paused with %d minutes left. Resolve the demand, then resume or abandon the job. No reward has been paid.", job.Title, w.SuspendedJob.Remaining), "story")
	} else {
		w.RememberArrangement(job, "interrupted")
		w.Log("An interrupted arrangement", "Danger stopped the operation. No reward was paid.", "story")
	}
}

func (w *World) OfferResume() {
	if w.SuspendedJob == nil || !w.Player.Alive || w.Event != nil {
		return
	}
	s := w.SuspendedJob
	w.Event = &Scene{ID: ID(), Kind: "resume_job", Source: "authored", Minute: w.Minute, Speaker: s.Scene.Speaker,
		Title: "Unfinished business", Body: fmt.Sprintf("The demand is settled. %s has %d minutes of work remaining. Its agreed terms have not changed.", s.Scene.Title, s.Remaining),
		Choices: []Choice{{ID: "resume", Label: "Resume the arrangement", Detail: fmt.Sprintf("%d minutes remaining · $%d reward · +%d respect · +%d heat. Police and further danger can still intervene.", s.Remaining, s.Scene.Effect.Reward, s.Scene.Effect.Respect, s.Scene.Effect.Heat)},
			{ID: "abandon", Label: "Abandon the unfinished work", Detail: "No additional time or reward. You can attend to other matters."}}}
}

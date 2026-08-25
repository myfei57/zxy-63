package ns

import "fmt"

type Step struct {
	Stage  Stage
	Action string
}

func (p Pipeline) Steps() []Step {
	actions := map[Stage]string{
		StageIntake:    "collect flow",
		StageCoag:      "dose coagulant",
		StageChlor:     "dose chlorine",
		StageFilter:    "filter water",
		StageClearwell: "store clear water",
		StageQuota:     "meter chemicals",
		StageAudit:     "record audit",
	}
	steps := make([]Step, 0, len(p.Stages))
	for _, stage := range p.Stages {
		steps = append(steps, Step{Stage: stage, Action: actions[stage]})
	}
	return steps
}

func (p Pipeline) Describe() string {
	return fmt.Sprintf("pipeline %s has %d stages", p.Name, len(p.Stages))
}

func (p Pipeline) Last() Stage {
	if len(p.Stages) == 0 {
		return ""
	}
	return p.Stages[len(p.Stages)-1]
}

func (p Pipeline) Count() int {
	return len(p.Stages)
}

func (p Pipeline) Contains(stage Stage) bool {
	return p.IndexOf(stage) >= 0
}

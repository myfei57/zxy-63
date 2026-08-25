package ns

type Stage string

const (
	StageIntake    Stage = "intake"
	StageCoag      Stage = "coag"
	StageChlor     Stage = "chlor"
	StageFilter    Stage = "filter"
	StageBackwash  Stage = "backwash"
	StageTurbidity Stage = "turbidity"
	StageClearwell Stage = "clearwell"
	StageQuota     Stage = "quota"
	StageAudit     Stage = "audit"
)

type Pipeline struct {
	Name   string
	Stages []Stage
}

func TreatmentLine() Pipeline {
	return Pipeline{
		Name:   "treatment",
		Stages: []Stage{StageIntake, StageCoag, StageChlor, StageFilter, StageClearwell, StageQuota, StageAudit},
	}
}

func (p Pipeline) IndexOf(stage Stage) int {
	for i, candidate := range p.Stages {
		if candidate == stage {
			return i
		}
	}
	return -1
}

func (p Pipeline) Before(a, b Stage) bool {
	ai := p.IndexOf(a)
	bi := p.IndexOf(b)
	return ai >= 0 && bi >= 0 && ai < bi
}

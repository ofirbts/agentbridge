package observability

type StepLog struct {
	RunID      string `json:"run_id"`
	Step       string `json:"step"`
	Status     string `json:"status"`
	Retries    int    `json:"retries"`
	DurationMS int64  `json:"duration_ms"`
	Provider   string `json:"provider"`
	Query      string `json:"query"`
	Message    string `json:"message,omitempty"`
}

type Tracer struct {
	logs []StepLog
}

func NewTracer() *Tracer {
	return &Tracer{logs: []StepLog{}}
}

func (t *Tracer) StartRun(runID string) *Recorder {
	return &Recorder{
		runID:  runID,
		tracer: t,
	}
}

func (t *Tracer) Logs() []StepLog {
	out := make([]StepLog, len(t.logs))
	copy(out, t.logs)
	return out
}
